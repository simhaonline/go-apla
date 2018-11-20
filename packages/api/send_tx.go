// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/http"

	"github.com/GenesisKernel/go-genesis/packages/block"
	"github.com/GenesisKernel/go-genesis/packages/conf/syspar"
	"github.com/GenesisKernel/go-genesis/packages/consts"
	"github.com/GenesisKernel/go-genesis/packages/converter"
	"github.com/GenesisKernel/go-genesis/packages/model"
	"github.com/GenesisKernel/go-genesis/packages/transaction"
	"github.com/GenesisKernel/go-genesis/packages/utils/tx"

	log "github.com/sirupsen/logrus"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type sendTxResult struct {
	Hashes map[string]string `json:"hashes"`
}

func getTxData(r *http.Request, key string) ([]byte, error) {
	logger := getLogger(r)

	file, _, err := r.FormFile(key)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Error("request.FormFile")
		return nil, err
	}
	defer file.Close()

	var txData []byte
	if txData, err = ioutil.ReadAll(file); err != nil {
		logger.WithFields(log.Fields{"type": consts.IOError, "error": err}).Error("reading multipart file")
		return nil, err
	}

	return txData, nil
}

func sendTxHandler(w http.ResponseWriter, r *http.Request) {
	client := getClient(r)

	if block.IsKeyBanned(client.KeyID) {
		errorResponse(w, errBannded.Errorf(block.BannedTill(client.KeyID)))
		return
	}

	err := r.ParseMultipartForm(multipartBuf)
	if err != nil {
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	result := &sendTxResult{Hashes: make(map[string]string)}
	for key := range r.MultipartForm.File {
		txData, err := getTxData(r, key)
		if err != nil {
			errorResponse(w, err)
			return
		}

		hash, err := txHandler(r, txData)
		if err != nil {
			errorResponse(w, err)
			return
		}
		result.Hashes[key] = hash
	}

	for key := range r.Form {
		txData, err := hex.DecodeString(r.FormValue(key))
		if err != nil {
			errorResponse(w, err)
			return
		}

		hash, err := txHandler(r, txData)
		if err != nil {
			errorResponse(w, err)
			return
		}
		result.Hashes[key] = hash
	}

	jsonResponse(w, result)
}

type contractResult struct {
	Hash string `json:"hash"`
	// These fields are used for VDE
	Message *txstatusError `json:"errmsg,omitempty"`
	Result  string         `json:"result,omitempty"`
}

func txHandler(r *http.Request, txData []byte) (string, error) {
	client := getClient(r)
	logger := getLogger(r)

	if int64(len(txData)) > syspar.GetMaxTxSize() {
		logger.WithFields(log.Fields{"type": consts.ParameterExceeded, "max_size": syspar.GetMaxTxSize(), "size": len(txData)}).Error("transaction size exceeds max size")
		block.BadTxForBan(client.KeyID)
		return "", errLimitTxSize.Errorf(len(txData))
	}

	rtx := &transaction.RawTransaction{}
	if err := rtx.Unmarshall(bytes.NewBuffer(txData)); err != nil {
		return "", err
	}
	smartTx := tx.SmartContract{}
	if err := msgpack.Unmarshal(rtx.Payload(), &smartTx); err != nil {
		return "", err
	}
	if smartTx.Header.KeyID != client.KeyID {
		return "", errDiffKey
	}
	if err := model.SendTx(rtx, client.KeyID); err != nil {
		logger.WithFields(log.Fields{"type": consts.DBError, "error": err}).Error("sending tx")
		return "", err
	}

	return string(converter.BinToHex(rtx.Hash())), nil
}
