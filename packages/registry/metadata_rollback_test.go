package registry

import (
	"testing"

	"encoding/json"

	"fmt"

	"strconv"

	"github.com/GenesisKernel/go-genesis/packages/storage/kv"
	"github.com/GenesisKernel/go-genesis/packages/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/yddmat/memdb"
)

type teststruct struct {
	Key    int
	Value1 string
	Value2 []byte
}

func TestMetadataRollbackSaveState(t *testing.T) {
	txMock := &kv.MockTransaction{}
	mr := metadataRollback{tx: txMock, txCounter: make(map[string]uint64)}

	registry := &types.Registry{
		Name:      "keys",
		Ecosystem: &types.Ecosystem{ID: 10},
	}

	block, tx := []byte("123"), []byte("321")

	s := state{Counter: 1, RegistryName: registry.Name, Ecosystem: registry.Ecosystem.ID, Key: "1"}
	jstate, err := json.Marshal(s)
	require.Nil(t, err)
	txMock.On("Set", fmt.Sprintf(writePrefix, string(block), 1, string(tx)), string(jstate)).Return(nil)
	require.Nil(t, mr.saveState(block, tx, registry, "1", ""))
	require.Equal(t, mr.txCounter[string(block)], uint64(1))

	structValue := teststruct{
		Key:    666,
		Value1: "stringvalue",
		Value2: make([]byte, 20),
	}
	jsonValue, err := json.Marshal(structValue)
	require.Nil(t, err)
	s = state{Counter: 2, RegistryName: registry.Name, Ecosystem: registry.Ecosystem.ID, Value: string(jsonValue), Key: "2"}
	jstate, err = json.Marshal(s)
	require.Nil(t, err)
	txMock.On("Set", fmt.Sprintf(writePrefix, string(block), 2, string(tx)), string(jstate)).Return(nil)
	require.Nil(t, mr.saveState(block, tx, registry, "2", string(jsonValue)))
	require.Equal(t, mr.txCounter[string(block)], uint64(2))

	s = state{Counter: 3, RegistryName: registry.Name, Ecosystem: registry.Ecosystem.ID, Value: "", Key: "3"}
	jstate, err = json.Marshal(s)
	require.Nil(t, err)
	txMock.On("Set", fmt.Sprintf(writePrefix, string(block), 3, string(tx)), string(jstate)).Return(errors.New("testerr"))
	require.Error(t, mr.saveState(block, tx, registry, "3", ""))
	require.Equal(t, mr.txCounter[string(block)], uint64(2))
}

func TestMetadataRollbackSaveRollback(t *testing.T) {
	mDb, err := memdb.OpenDB("", false)
	require.Nil(t, err)
	db := kv.DatabaseAdapter{Database: *mDb}

	dbTx := db.Begin(true)
	mr := metadataRollback{tx: dbTx, txCounter: make(map[string]uint64)}

	registry := &types.Registry{
		Name:      "keys",
		Ecosystem: &types.Ecosystem{ID: 5},
	}

	block := []byte("123")

	for key := range make([]int, 20) {
		// Emulating new value in database
		dbTx.Set(fmt.Sprintf(keyConvention, registry.Name, registry.Ecosystem.ID, strconv.Itoa(key)), "{\"result\":\"blah\"")

		tx := []byte(strconv.Itoa(key))
		tx = append(tx, []byte("blah")...)

		value := teststruct{
			Key:    key,
			Value1: "stringvalue" + strconv.Itoa(key),
			Value2: make([]byte, 20),
		}

		jsonValue, err := json.Marshal(value)
		require.Nil(t, err)
		// Save "old" state of record
		require.Nil(t, mr.saveState(block, tx, registry, strconv.Itoa(key), string(jsonValue)))
	}
	require.Nil(t, dbTx.Commit())

	dbTx = db.Begin(false)

	// We need to check that all previous states was saved to db
	for key := range make([]int, 20) {
		tx := []byte(strconv.Itoa(key))
		tx = append(tx, []byte("blah")...)

		_, err := dbTx.Get(fmt.Sprintf(writePrefix, string(block), key+1, string(tx)))
		require.Nil(t, err)
	}
	require.Nil(t, dbTx.Commit())

	dbTx = db.Begin(true)
	require.Nil(t, err)

	dbTx.AddIndex(types.Index{Name: "rollback", Registry: &types.Registry{Name: "rollback_tx", Ecosystem: &types.Ecosystem{ID: 5}}, SortFn: func(a, b string) bool {
		return true
	}})

	mr = metadataRollback{tx: dbTx, txCounter: make(map[string]uint64)}
	require.Nil(t, mr.rollbackState(block))

	// We are checking that all values are now at the previous state
	for key := range make([]int, 20) {
		// Emulating new value in database
		value, err := dbTx.Get(fmt.Sprintf(keyConvention, registry.Name, registry.Ecosystem.ID, strconv.Itoa(key)))
		require.Nil(t, err)

		got := teststruct{}
		json.Unmarshal([]byte(value), &got)
		require.Equal(t, teststruct{
			Key:    key,
			Value1: "stringvalue" + strconv.Itoa(key),
			Value2: make([]byte, 20),
		}, got)
	}
}
