package main

import (
	"math/rand"
	"testing"

	"singula.finance/node/chain"

	"singula.finance/node/block"

	"singula.finance/netio/crypto"
)

//basic chain functions
func TestChainsetup(t *testing.T) {

	mgr := chain.CreateManager()
	mgr.InitAccounts()
	//ra := mgr.RandomAccount()
	//log.Println(ra)

	//Genesis_Account := block.AccountFromString(chain.Treasury_Address)
	randNonce := rand.Intn(100)
	r := crypto.RandomPublicKey()
	address_r := crypto.Address(r)
	//r_account := block.AccountFromString(address_r)
	amount := 10

	someTx := block.Tx{Nonce: randNonce, Sender: chain.Treasury_Address, Receiver: address_r, Amount: amount}
	//log.Println(someTx)

	b := block.Block{}
	b.Txs = []block.Tx{}
	b.Txs = append(b.Txs, someTx)
	mgr.ApplyBlock(b)

	if mgr.State.Accounts[address_r] != 10 {
		t.Error("wrong amount")
	}

	//chain.InitAccounts()
}

func TestTxf(t *testing.T) {

	//create genesis block

	//test LastBlock

}
