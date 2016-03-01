package test

import (
	"bytes"
	"testing"

	. "github.com/tendermint/go-common"
	"github.com/porkchop/merklepatricia/app"
	patricia "github.com/porkchop/merklepatricia/client"
	"github.com/tendermint/tmsp/server"
)

func TestClient(t *testing.T) {

	addr := "tcp://127.0.0.1:46659"

	// Start the listener
	mApp := app.NewMerklePatriciaApp()
	s, err := server.NewServer(addr, mApp)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	defer s.Stop()

	// Create client
	cli, err := patricia.NewClient(addr)
	if err != nil {
		t.Fatal(err.Error())
		return
	}
	defer cli.Stop()

	// Empty
	commit(t, cli, "")
	get(t, cli, "foo", "", "")
	get(t, cli, "bar", "", "")
	// Set foo=FOO
	set(t, cli, "foo", "FOO")
	commit(t, cli, "E25690726078298180EB98AD2734F69FB8EFEECAA636FB58E3E24ADED1A40C1F")
	get(t, cli, "foo", "FOO", "")
	get(t, cli, "foa", "", "")
	get(t, cli, "foz", "", "")
	rem(t, cli, "foo")
	// Empty
	get(t, cli, "foo", "", "")
	commit(t, cli, "")
	// Set foo1, foo2, foo3...
	set(t, cli, "foo1", "1")
	set(t, cli, "foo2", "2")
	set(t, cli, "foo3", "3")
	set(t, cli, "foo1", "4")
	get(t, cli, "foo1", "4", "")
	get(t, cli, "foo2", "2", "")
	get(t, cli, "foo3", "3", "")
	rem(t, cli, "foo3")
	rem(t, cli, "foo2")
	rem(t, cli, "foo1")
	// Empty
	commit(t, cli, "")

}

func get(t *testing.T, cli *patricia.Client, key string, value string, err string) {
	_value, _err := cli.GetSync([]byte(key))
	if !bytes.Equal([]byte(value), _value) {
		t.Errorf("Expected value 0x%X (%v) but got 0x%X", []byte(value), value, _value)
	}
	if _err == nil {
		if err != "" {
			t.Errorf("Expected error %v but got no error", err)
		}
	} else {
		if err == "" {
			t.Errorf("Expected no error but got error %v", _err.Error())
		}
	}
}

func set(t *testing.T, cli *patricia.Client, key string, value string) {
	cli.SetSync([]byte(key), []byte(value))
}

func rem(t *testing.T, cli *patricia.Client, key string) {
	cli.RemSync([]byte(key))
}

func commit(t *testing.T, cli *patricia.Client, hash string) {
	_hash, _, err := cli.CommitSync()
	if err != nil {
		t.Error("Unexpected error getting hash", err.Error())
	}
	if hash != Fmt("%X", _hash) {
		t.Errorf("Expected hash 0x%v but got 0x%X", hash, _hash)
	}
}
