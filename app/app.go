package app

import (
	"fmt"
	. "github.com/tendermint/go-common"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-wire"
	tmsp "github.com/tendermint/tmsp/types"
)

type MerklePatriciaApp struct {
	trie trie.Trie
}

func NewMerklePatriciaApp() *MerklePatriciaApp {
	db, _ := ethdb.NewMemDatabase()
	_trie, _ := trie.New(
		common.Hash{},
		db,
	)
	return &MerklePatriciaApp{trie: *_trie}
}

func size(_trie *trie.Trie) (size int) {
	it := trie.NewIterator(_trie)
	for it.Next() {
	   size++
	}
	return
}

func (app *MerklePatriciaApp) Info() string {
	return Fmt("size:%v", size(&app.trie))
}

func (app *MerklePatriciaApp) SetOption(key string, value string) (log string) {
	return "No options are supported yet"
}

func (app *MerklePatriciaApp) AppendTx(tx []byte) (code tmsp.CodeType, result []byte, log string) {
	if len(tx) == 0 {
		return tmsp.CodeType_EncodingError, nil, "Tx length cannot be zero"
	}
	typeByte := tx[0]
	tx = tx[1:]
	switch typeByte {
	case 0x01: // Set
		key, n, err := wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting key: %v", err.Error())
		}
		tx = tx[n:]
		value, n, err := wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting value: %v", err.Error())
		}
		tx = tx[n:]
		if len(tx) != 0 {
			return tmsp.CodeType_EncodingError, nil, Fmt("Got bytes left over")
		}
		app.trie.Update(key, value)
		fmt.Println("SET", Fmt("%X", key), Fmt("%X", value))
	case 0x02: // Remove
		key, n, err := wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting key: %v", err.Error())
		}
		tx = tx[n:]
		if len(tx) != 0 {
			return tmsp.CodeType_EncodingError, nil, Fmt("Got bytes left over")
		}
		app.trie.Delete(key)
	default:
		return tmsp.CodeType_UnknownRequest, nil, Fmt("Unexpected type byte %X", typeByte)
	}
	return tmsp.CodeType_OK, nil, ""
}

func (app *MerklePatriciaApp) CheckTx(tx []byte) (code tmsp.CodeType, result []byte, log string) {
	if len(tx) == 0 {
		return tmsp.CodeType_EncodingError, nil, "Tx length cannot be zero"
	}
	typeByte := tx[0]
	tx = tx[1:]
	switch typeByte {
	case 0x01: // Set
		_, n, err := wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting key: %v", err.Error())
		}
		tx = tx[n:]
		_, n, err = wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting value: %v", err.Error())
		}
		tx = tx[n:]
		if len(tx) != 0 {
			return tmsp.CodeType_EncodingError, nil, Fmt("Got bytes left over")
		}
		//app.trie.Update(key, value)
	case 0x02: // Remove
		_, n, err := wire.GetByteSlice(tx)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting key: %v", err.Error())
		}
		tx = tx[n:]
		if len(tx) != 0 {
			return tmsp.CodeType_EncodingError, nil, Fmt("Got bytes left over")
		}
		//app.trie.Delete(key)
	default:
		return tmsp.CodeType_UnknownRequest, nil, Fmt("Unexpected type byte %X", typeByte)
	}
	return tmsp.CodeType_OK, nil, ""
}

func (app *MerklePatriciaApp) Commit() (hash []byte, log string) {
	if size(&app.trie) == 0 {
		return nil, "Empty hash for empty tree"
	}
	hash = app.trie.Hash().Bytes()
	return hash, ""
}

func (app *MerklePatriciaApp) Query(query []byte) (code tmsp.CodeType, result []byte, log string) {
	if len(query) == 0 {
		return tmsp.CodeType_OK, nil, "Query length cannot be zero"
	}
	typeByte := query[0]
	query = query[1:]
	switch typeByte {
	case 0x01: // Get
		key, n, err := wire.GetByteSlice(query)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error getting key: %v", err.Error())
		}
		query = query[n:]
		if len(query) != 0 {
			return tmsp.CodeType_EncodingError, nil, Fmt("Got bytes left over")
		}
		value := app.trie.Get(key)
		res := make([]byte, wire.ByteSliceSize(value))
		buf := res
		n, err = wire.PutByteSlice(buf, value)
		if err != nil {
			return tmsp.CodeType_EncodingError, nil, Fmt("Error putting value: %v", err.Error())
		}
		buf = buf[n:]
		return tmsp.CodeType_OK, res, ""
	default:
		return tmsp.CodeType_UnknownRequest, nil, Fmt("Unexpected type byte %X", typeByte)
	}
}
