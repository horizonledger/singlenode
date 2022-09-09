package chain

import (
	"encoding/hex"
	"log"
	"testing"

	"singula.finance/node/block"

	"golang.org/x/exp/maps"
	"singula.finance/netio/crypto"
)

func TestBalance(t *testing.T) {

	mgr := CreateManager()
	mgr.InitAccounts()

	ks := maps.Keys(mgr.State.Accounts)

	if mgr.State.Accounts[ks[0]] == 0 {
		t.Error("0 balance")
	}

	genBlock := MakeGenesisBlock()
	if genBlock.Height != 0 {
		t.Error("wrong height")
	}

	mgr.ApplyAppendBlock(genBlock)

	if mgr.BlockHeight() != 1 {
		t.Error("wrong height")
	}

	// if reply.MessageType != netio.REP {
	// 	t.Error("balance msg", reply_msg)
	// }

}

func TestTx(t *testing.T) {
	//initialize blockchain
	mgr := CreateManager()
	mgr.InitAccounts()
	genBlock := MakeGenesisBlock()
	mgr.ApplyAppendBlock(genBlock)
	EmptyPool(&mgr)
	//mgr.WriteChain()
	//create accounts
	keypair_sender := crypto.PairFromSecret("sender")
	pubkey_sender := crypto.PubKeyToHex(keypair_sender.PubKey)
	addr_sender := crypto.Address(pubkey_sender)
	mgr.SetAccount(addr_sender, 100)
	if mgr.State.Accounts[addr_sender] != 100 {
		t.Error("wrong balance sender")
	}

	log.Println("initial balance sender = ", mgr.State.Accounts[addr_sender])
	keypair_receiver := crypto.PairFromSecret("receiver")
	addr_receiver := crypto.Address(crypto.PubKeyToHex(keypair_receiver.PubKey))
	mgr.SetAccount(addr_receiver, 0)
	if mgr.State.Accounts[addr_receiver] != 0 {
		t.Error("wrong balance sender")
	}

	log.Println("initial balance receiver = ", mgr.State.Accounts[addr_receiver])
	// create tx
	tx := block.Tx{Nonce: 1, Amount: 10, Sender: addr_sender, Receiver: addr_receiver, SenderPubkey: pubkey_sender}
	// sign tx
	signature := block.SignTx(tx, keypair_sender.PrivKey)
	tx.Signature = hex.EncodeToString(signature.Serialize())
	// validate tx
	// add tx to block
	HandleTx(&mgr, tx) //verify tx validity, add tx to Tx_pool and broadcast tx
	//[ ] check block calculations

	if len(mgr.Tx_pool) != 1 {
		t.Error("tx pool not full")
	}
	if mgr.BlockHeight() != 1 {
		t.Error("wrong height ", mgr.BlockHeight())
	}
	// append block
	MakeBlock(&mgr) //create a block with tx (Apply tx, Append Block, empty Tx_pool)
	if mgr.BlockHeight() != 2 {
		t.Error("wrong height after apply block", mgr.BlockHeight())
	}
	//after tx is confirmed, balance should have changed
	//log.Println("final balance sender = ", mgr.State.Accounts[addr_sender])
	if !(mgr.State.Accounts[tx.Sender] == 90) {
		t.Error("sender wrong balance")
	}
	//log.Println("final balance receiver = ", mgr.State.Accounts[addr_receiver])
	if !(mgr.State.Accounts[tx.Receiver] == 10) {
		t.Error("receiver wrong balance")
	}
	//[ ] sign block
	//[ ] check signature block
	// no function for this in programs yet!

}
