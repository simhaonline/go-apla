// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package daemons

import (
	"context"
	"io/ioutil"
	"sync/atomic"
	"time"

	"github.com/AplaProject/go-apla/packages/block"
	"github.com/AplaProject/go-apla/packages/network"
	"github.com/AplaProject/go-apla/packages/network/tcpclient"

	"github.com/AplaProject/go-apla/packages/conf"
	"github.com/AplaProject/go-apla/packages/conf/syspar"
	"github.com/AplaProject/go-apla/packages/consts"
	"github.com/AplaProject/go-apla/packages/crypto"
	"github.com/AplaProject/go-apla/packages/model"
	"github.com/AplaProject/go-apla/packages/rollback"
	"github.com/AplaProject/go-apla/packages/service"
	"github.com/AplaProject/go-apla/packages/transaction"
	"github.com/AplaProject/go-apla/packages/utils"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNodesUnavailable is returned when all nodes is unavailable
	ErrNodesUnavailable = errors.New("All nodes unavailable")
)

// BlocksCollection collects and parses blocks
func BlocksCollection(ctx context.Context, d *daemon) error {
	if ctx.Err() != nil {
		d.logger.WithFields(log.Fields{"type": consts.ContextError, "error": ctx.Err()}).Error("context error")
		return ctx.Err()
	}

	return blocksCollection(ctx, d)
}

func InitialLoad(logger *log.Entry) error {

	// check for initial load
	toLoad, err := needLoad(logger)
	if err != nil {
		return err
	}

	if toLoad {
		logger.Debug("start first block loading")

		if err := firstLoad(logger); err != nil {
			return err
		}

		model.UpdateSchema()
	}

	return nil
}

var bcOnRun uint32

func blocksCollection(ctx context.Context, d *daemon) (err error) {
	if !atomic.CompareAndSwapUint32(&bcOnRun, 0, 1) {
		return nil
	}
	defer func() {
		atomic.StoreUint32(&bcOnRun, 0)
	}()
	host, maxBlockID, err := getHostWithMaxID(ctx, d.logger)
	if err != nil {
		d.logger.WithFields(log.Fields{"error": err}).Warn("on checking best host")
		return err
	}

	infoBlock := &model.InfoBlock{}
	found, err := infoBlock.Get()
	if err != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting cur blockID")
		return err
	}
	if !found {
		log.WithFields(log.Fields{"type": consts.NotFound, "error": err}).Error("Info block not found")
		return errors.New("Info block not found")
	}

	if infoBlock.BlockID >= maxBlockID {
		log.WithFields(log.Fields{"blockID": infoBlock.BlockID, "maxBlockID": maxBlockID}).Debug("Max block is already in the host")
		return nil
	}

	DBLock()
	defer func() {
		service.NodeDoneUpdatingBlockchain()
		DBUnlock()
	}()

	// update our chain till maxBlockID from the host
	return UpdateChain(ctx, d, host, maxBlockID)
}

// UpdateChain load from host all blocks from our last block to maxBlockID
func UpdateChain(ctx context.Context, d *daemon, host string, maxBlockID int64) (err error) {
	// get current block id from our blockchain
	curBlock := &model.InfoBlock{}
	if _, err = curBlock.Get(); err != nil {
		d.logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("Getting info block")
		return err
	}

	if ctx.Err() != nil {
		d.logger.WithFields(log.Fields{"type": consts.ContextError, "error": ctx.Err()}).Error("context error")
		return ctx.Err()
	}

	var lastBlockID, lastBlockTime int64
	defer func() {
		if err != nil {
			banNode(host, lastBlockID, lastBlockTime, err)
		}
	}()

	playRawBlock := func(rb []byte) error {
		bl, err := block.ProcessBlockWherePrevFromBlockchainTable(rb, true)
		if err != nil {
			d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("processing block")
			return err
		}

		lastBlockID = bl.Header.BlockID
		lastBlockTime = bl.Header.Time

		// hash compare could be failed in the case of fork
		hashMatched, errCheck := bl.CheckHash()
		if errCheck != nil {
			d.logger.WithFields(log.Fields{"error": errCheck, "type": consts.BlockError}).Error("checking block hash")
		}

		if !hashMatched {
			transaction.CleanCache()

			var replaceCount int64 = 1
			if errCheck == block.ErrIncorrectRollbackHash {
				replaceCount++
			}

			//it should be fork, replace our previous blocks to ones from the host
			err = ReplaceBlocksFromHost(ctx, host, bl.Header.BlockID-1, replaceCount)
			if err != nil {
				d.logger.WithFields(log.Fields{"error": err, "type": consts.ParserError}).Error("replacing blocks")
				return err
			}
		}

		bl.PrevHeader, err = block.GetBlockDataFromBlockChain(bl.Header.BlockID - 1)
		if err != nil {
			return errors.Wrapf(err, "can't get block %d", bl.Header.BlockID-1)
		}

		if err = bl.Check(); err != nil {
			return err
		}

		return bl.PlaySafe()
	}

	var count int
	st := time.Now()

	d.logger.WithFields(log.Fields{"min_block": curBlock.BlockID, "max_block": maxBlockID, "count": maxBlockID - curBlock.BlockID}).Info("starting downloading blocks")
	for blockID := curBlock.BlockID + 1; blockID <= maxBlockID; blockID += int64(network.BlocksPerRequest) {

		if loopErr := func() error {
			ctxDone, cancel := context.WithCancel(ctx)
			defer func() {
				cancel()
				d.logger.WithFields(log.Fields{"count": count, "time": time.Since(st).String()}).Info("blocks downloaded")
			}()

			rawBlocksChan, err := tcpclient.GetBlocksBodies(ctxDone, host, blockID, false)
			if err != nil {
				d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("getting block body")
				return err
			}

			for rawBlock := range rawBlocksChan {
				if err = playRawBlock(rawBlock); err != nil {
					d.logger.WithFields(log.Fields{"error": err, "type": consts.BlockError}).Error("playing raw block")
					return err
				}
				count++
			}

			return nil
		}(); loopErr != nil {
			return loopErr
		}
	}
	return nil
}

// init first block from file or from embedded value
func loadFirstBlock(logger *log.Entry) error {
	newBlock, err := ioutil.ReadFile(conf.Config.FirstBlockPath)
	if err != nil {
		logger.WithFields(log.Fields{
			"type": consts.IOError, "error": err, "path": conf.Config.FirstBlockPath,
		}).Error("reading first block from file")
	}

	if err = block.InsertBlockWOForks(newBlock, false, true); err != nil {
		logger.WithFields(log.Fields{"type": consts.ParserError, "error": err}).Error("inserting new block")
		return err
	}

	return nil
}

func firstLoad(logger *log.Entry) error {
	DBLock()
	defer DBUnlock()

	return loadFirstBlock(logger)
}

func needLoad(logger *log.Entry) (bool, error) {
	infoBlock := &model.InfoBlock{}
	_, err := infoBlock.Get()
	if err != nil {
		logger.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("getting info block")
		return false, err
	}
	// we have empty blockchain, we need to load blockchain from file or other source
	if infoBlock.BlockID == 0 {
		logger.Debug("blockchain should be loaded")
		return true, nil
	}
	return false, nil
}

func banNode(host string, blockID, blockTime int64, err error) {
	// disable ban node
	return

	if err == nil || !utils.IsBanError(err) {
		return
	}

	reason := err.Error()
	log.WithFields(log.Fields{"reason": reason, "host": host, "block_id": blockID, "block_time": blockTime}).Info("ban node")

	n, err := syspar.GetNodeByHost(host)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("getting node by host")
		return
	}

	err = service.GetNodesBanService().RegisterBadBlock(n, blockID, blockTime, reason)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "node": n.KeyID, "block": blockID}).Error("registering bad block from node")
	}
}

// GetHostWithMaxID returns host with maxBlockID
func getHostWithMaxID(ctx context.Context, logger *log.Entry) (host string, maxBlockID int64, err error) {

	nbs := service.GetNodesBanService()
	hosts, err := nbs.FilterBannedHosts(syspar.GetRemoteHosts())
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("on filtering banned hosts")
	}

	host, maxBlockID, err = tcpclient.HostWithMaxBlock(ctx, hosts)
	if len(hosts) == 0 || err == tcpclient.ErrNodesUnavailable {
		hosts = conf.GetNodesAddr()
		return tcpclient.HostWithMaxBlock(ctx, hosts)
	}

	return
}

// ReplaceBlocksFromHost replaces blockchain received from the host.
// Number (replaceCount) of blocks starting from blockID will be re-played.
func ReplaceBlocksFromHost(ctx context.Context, host string, blockID, replaceCount int64) error {
	blocks, err := getBlocks(ctx, host, blockID, replaceCount)
	if err != nil {
		return err
	}
	transaction.CleanCache()

	// mark all transaction as unverified
	_, err = model.MarkVerifiedAndNotUsedTransactionsUnverified()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"type":  consts.DBError,
		}).Error("marking verified and not used transactions unverified")
		return utils.ErrInfo(err)
	}

	// get starting blockID from slice of blocks
	if len(blocks) > 0 {
		blockID = blocks[len(blocks)-1].Header.BlockID
	}

	// we have the slice of blocks for applying
	// first of all we should rollback old blocks
	block := &model.Block{}
	myRollbackBlocks, err := block.GetBlocksFrom(blockID-1, "desc", 0)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("getting rollback blocks from blockID")
		return utils.ErrInfo(err)
	}
	for _, block := range myRollbackBlocks {
		err := rollback.RollbackBlock(block.Data)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}

	return processBlocks(blocks)
}

func getBlocks(ctx context.Context, host string, blockID, minCount int64) ([]*block.Block, error) {
	rollback := syspar.GetRbBlocks1()
	blocks := make([]*block.Block, 0)
	nextBlockID := blockID

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// load the block bodies from the host
	blocksCh, err := tcpclient.GetBlocksBodies(ctx, host, blockID, true)
	if err != nil {
		return nil, utils.WithBan(errors.Wrapf(err, "Getting bodies of blocks by id %d", blockID))
	}

	for binaryBlock := range blocksCh {
		if blockID < 2 {
			break
		}

		// if the limit of blocks received from the node was exaggerated
		if len(blocks) >= int(rollback) {
			break
		}

		bl, err := block.ProcessBlockWherePrevFromBlockchainTable(binaryBlock, true)
		if err != nil {
			return nil, err
		}

		if bl.Header.BlockID != nextBlockID {
			log.WithFields(log.Fields{"header_block_id": bl.Header.BlockID, "block_id": blockID, "type": consts.InvalidObject}).Error("block ids does not match")
			return nil, utils.WithBan(errors.New("bad block_data['block_id']"))
		}

		// the public key of the one who has generated this block
		nodePublicKey, err := syspar.GetNodePublicKeyByPosition(bl.Header.NodePosition)
		if err != nil {
			log.WithFields(log.Fields{"header_block_id": bl.Header.BlockID, "block_id": blockID, "type": consts.InvalidObject}).Error("block ids does not match")
			return nil, utils.ErrInfo(err)
		}

		// save the block
		blocks = append(blocks, bl)

		// check the signature
		_, okSignErr := utils.CheckSign([][]byte{nodePublicKey},
			[]byte(bl.Header.ForSign(bl.PrevHeader, bl.MrklRoot)),
			bl.Header.Sign, true)
		if okSignErr == nil && len(blocks) >= int(minCount) {
			break
		}

		nextBlockID--
	}

	return blocks, nil
}

func processBlocks(blocks []*block.Block) error {
	dbTransaction, err := model.StartTransaction()
	if err != nil {
		log.WithFields(log.Fields{"error": err, "type": consts.DBError}).Error("starting transaction")
		return utils.ErrInfo(err)
	}

	// go through new blocks from the smallest block_id to the largest block_id
	prevBlocks := make(map[int64]*block.Block, 0)

	for i := len(blocks) - 1; i >= 0; i-- {
		b := blocks[i]

		if prevBlocks[b.Header.BlockID-1] != nil {
			b.PrevHeader.Hash = prevBlocks[b.Header.BlockID-1].Header.Hash
			b.PrevHeader.RollbacksHash = prevBlocks[b.Header.BlockID-1].Header.RollbacksHash
			b.PrevHeader.Time = prevBlocks[b.Header.BlockID-1].Header.Time
			b.PrevHeader.BlockID = prevBlocks[b.Header.BlockID-1].Header.BlockID
			b.PrevHeader.EcosystemID = prevBlocks[b.Header.BlockID-1].Header.EcosystemID
			b.PrevHeader.KeyID = prevBlocks[b.Header.BlockID-1].Header.KeyID
			b.PrevHeader.NodePosition = prevBlocks[b.Header.BlockID-1].Header.NodePosition
		}

		hash, err := crypto.DoubleHash([]byte(b.Header.ForSha(b.PrevHeader, b.MrklRoot)))
		if err != nil {
			log.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Fatal("double hashing block")
		}
		b.Header.Hash = hash

		if err := b.Check(); err != nil {
			dbTransaction.Rollback()
			return err
		}

		if err := b.Play(dbTransaction); err != nil {
			dbTransaction.Rollback()
			return utils.ErrInfo(err)
		}
		prevBlocks[b.Header.BlockID] = b

		// for last block we should update block info
		if i == 0 {
			err := block.UpdBlockInfo(dbTransaction, b)
			if err != nil {
				dbTransaction.Rollback()
				return utils.ErrInfo(err)
			}
		}
		if b.SysUpdate {
			if err := syspar.SysUpdate(dbTransaction); err != nil {
				log.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("updating syspar")
				return utils.ErrInfo(err)
			}
		}
	}

	// If all right we can delete old blockchain and write new
	for i := len(blocks) - 1; i >= 0; i-- {
		b := blocks[i]
		// insert new blocks into blockchain
		if err := block.InsertIntoBlockchain(dbTransaction, b); err != nil {
			dbTransaction.Rollback()
			return err
		}
	}

	return dbTransaction.Commit()
}
