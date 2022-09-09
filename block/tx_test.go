package block

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"

	"testing"

	"singula.finance/netio/crypto"
)

// import (
// 	"encoding/hex"
// 	"encoding/json"
// 	"io/ioutil"
// 	"os"
// 	"singula/node/chain"
// 	"testing"

// 	"singula.finance/netio"
// )

// basic block functions
func TestBasicAssign(t *testing.T) {
	var tx Tx
	tx = Tx{Nonce: 1}
	if tx.Nonce != 1 {
		t.Error("fail assign nonce")
	}
}

func TestTxJson(t *testing.T) {
	var tx Tx
	tx = Tx{Nonce: 1}
	txJson, _ := json.Marshal(tx)
	if txJson[0] != '{' {
		t.Error("start json")
	}
	i := len(txJson) - 1
	if txJson[i] != '}' {
		t.Error("end json")
	}

	var newtx Tx
	if err := json.Unmarshal(txJson, &newtx); err != nil {
		panic(err)
	}
	if newtx.Nonce != tx.Nonce {
		t.Error("json marshal failed")
	}
	if newtx.Sender != tx.Sender {
		t.Error("json marshal failed")
	}
}

func TestSignTx(t *testing.T) {
	//sign
	keypair := crypto.PairFromSecret("test")
	var tx Tx
	//s := block.AccountFromString("Pa033f6528cc1")
	s := "Pa033f6528cc1"
	r := s //TODO
	tx = Tx{Nonce: 0, Amount: 0, Sender: s, Receiver: r}

	signature := SignTx(tx, keypair.PrivKey)
	sighex := hex.EncodeToString(signature.Serialize())

	if sighex == "" {
		t.Error("hex empty")
	}
	tx.Signature = sighex
	tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

	//verify
	verified := VerifyTxSig(tx)

	if !verified {
		t.Error("verify tx fail")
	}

}

func TestSignTxBasic(t *testing.T) {

	keypair := crypto.PairFromSecret("test")
	pub := crypto.PubKeyToHex(keypair.PubKey)
	//account := Account{AccountKey: crypto.Address(pub)}

	randNonce := 0
	amount := 10

	addr := crypto.Address(crypto.PubKeyToHex(keypair.PubKey))

	//crypto.RandomPublicKey()

	//genkeypair := GenesisKeys()
	// addr := crypto.Address(crypto.PubKeyToHex(genkeypair.PubKey))
	// //Genesis_Account := AccountFromString(addr)

	// // //{"Nonce":0,"Amount":0,"Sender":{"AccountKey":"Pa033f6528cc1"},"Receiver":{"AccountKey":"Pa033f6528cc1"},"SenderPubkey":"","Signature":"","id":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0]}
	// r := crypto.Address(pub)
	tx := Tx{Nonce: randNonce, Amount: amount, Sender: addr, Receiver: addr, SenderPubkey: pub, Signature: ""}

	if tx.Amount != 10 {
		t.Error("wrong amount")
	}

	//tx = crypto.SignTxAdd(tx, keypair)

	// //log.Println(tx)

	//verified := crypto.VerifyTxSig(tx)

	// if !verified {
	// 	t.Error("verify tx fail")
	// }
}

func TestTxFile(t *testing.T) {
	//write tx.json

	keypair := crypto.PairFromSecret("test")

	pubk := crypto.PubKeyToHex(keypair.PubKey)
	addr := crypto.Address(pubk)

	if addr != "Pa033f6528cc1" {
		t.Error("address wrong ", addr)
	}

	keypair_recv := crypto.PairFromSecret("receive")
	addr_recv := crypto.Address(crypto.PubKeyToHex(keypair_recv.PubKey))

	tx := Tx{Nonce: 1, Amount: 10, Sender: addr, Receiver: addr_recv}
	signature := SignTx(tx, keypair.PrivKey)
	sighex := hex.EncodeToString(signature.Serialize())

	tx.Signature = sighex
	tx.SenderPubkey = crypto.PubKeyToHex(keypair.PubKey)

	if !(tx.Amount == 10) {
		t.Error("amount wrong ")
	}

	txJson, _ := json.Marshal(tx)

	ioutil.WriteFile("tx_test.json", []byte(txJson), 0644)

	dat, _ := ioutil.ReadFile("tx_test.json")

	os.Remove("tx_test.json")

	var newTx Tx

	if err := json.Unmarshal(dat, &newTx); err != nil {
		panic(err)
	}

	if !(newTx.Amount == 10) {
		t.Error("amount wrong ")
	}

	if newTx.Amount != tx.Amount {
		t.Error("amount not equal")
	}

	if newTx.Receiver != tx.Receiver {
		t.Error("Receiver not equal")
	}

	if newTx.Sender != tx.Sender {
		t.Error("Sender not equal")
	}

	//TODO more test for equal

	//log.Println(rtx.SenderPubkey)

	verified := VerifyTxSig(tx)

	if !verified {
		t.Error("verify tx fail")
	}

}

func TestSeralizeTx(t *testing.T) {

}

/* func TestTx(t *testing.T) {
	//initialize blockchain
	mgr := chain.CreateManager()
	mgr.InitAccounts()
	genBlock := chain.MakeGenesisBlock()
	mgr.ApplyBlock(genBlock)
	mgr.AppendBlock(genBlock)
	chain.EmptyPool(&mgr)
	mgr.WriteChain()
	//create addresses
	keypair_sender := crypto.PairFromSecret("sender")
	pubkey_sender := crypto.PubKeyToHex(keypair_sender.PubKey)
	addr_sender := crypto.Address(pubkey_sender)
	mgr.SetAccount(addr_sender, 100)
	keypair_receiver := crypto.PairFromSecret("receiver")
	addr_receiver := crypto.Address(crypto.PubKeyToHex(keypair_receiver.PubKey))
	mgr.SetAccount(addr_sender, 0)
	//[ ] create tx
	tx := Tx{Nonce: 1, Amount: 10, Sender: addr_sender, Receiver: addr_receiver, SenderPubkey: pubkey_sender}
	//[ ] sign tx
	signature := SignTx(tx, keypair_sender.PrivKey)
	tx.Signature = hex.EncodeToString(signature.Serialize())
	//[ ] validate tx
	chain.HandleTx(&mgr, tx) //verify tx validity, add tx to Tx_pool and broadcast tx
	chain.MakeBlock(&mgr)    //create a block with tx (Apply tx, Append Block, empty Tx_pool)
	//[ ] after tx is confirmed, balance should have changed
	if !(mgr.State.Accounts[tx.Sender] == 90) {
		t.Error("sender wrong balance")
	}
	if !(mgr.State.Accounts[tx.Sender] == 20) {
		t.Error("receiver wrong balance")
	}
	//[ ] add tx to block
	//[ ] check block calculations
	//[ ] append block
	//[ ] sign block
	//[ ] check signature block
} */

// func TestTxSign(t *testing.T) {
// 	kp := crypto.PairFromSecret("test")
// 	tx := Tx{TxType: "STX", Amount: 10, Sender: "Pa033f6528cc1", Receiver: "P7ba453f23337", Nonce: 0}
// 	signedtx := crypto.CreateSignedTx(tx, kp)
// 	signedtxjson, _ := json.Marshal(signedtx)
// 	if string(signedtxjson) != `{"txType":"STX","amount":10,"sender":"Pa033f6528cc1","receiver":"P7ba453f23337","nonce":0,"senderPubkey":"03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31","signature":"30450221009d3e5b449ad4870752e917906e379b82bcee234efb9fb8475541e2ca2066a431022053e7b5fa4316caefc5a0a039bdbb934e3f210d83f32097b614aae8ccb0bd6188"}` {
// 		t.Error("signedtxjson ", signedtxjson)
// 	}
// 	//ioutil.WriteFile("example.txps", []byte(signedtxjson), 0644)

// }

// func TestTxSignHandleMsg(t *testing.T) {

// 	kp := crypto.PairFromSecret("test")
// 	tx := Tx{TxType: "STX", Amount: 10, Sender: "Pa033f6528cc1", Receiver: "P7ba453f23337", Nonce: 0}
// 	signedtx := crypto.CreateSignedTx(tx, kp)
// 	signedtxjson, _ := json.Marshal(signedtx)
// 	txmsg := netio.Message{MessageType: netio.REQ, Command: netio.CMD_TX, Data: signedtxjson}
// 	txmsg_json, _ := json.Marshal(txmsg)
// 	if string(txmsg_json) != `{"MessageType":"REQ","Command":"TX","Data":{"txType":"STX","amount":10,"sender":"Pa033f6528cc1","receiver":"P7ba453f23337","nonce":0,"senderPubkey":"03dab2d148f103cd4761df382d993942808c1866a166f27cafba3289e228384a31","signature":"30450221009d3e5b449ad4870752e917906e379b82bcee234efb9fb8475541e2ca2066a431022053e7b5fa4316caefc5a0a039bdbb934e3f210d83f32097b614aae8ccb0bd6188"}}` {
// 		t.Error("tx msg", string(txmsg_json))
// 	}
// }
