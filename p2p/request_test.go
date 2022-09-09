package p2p

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"singula.finance/netio"
	"singula.finance/netio/crypto"
	"singula.finance/node/block"
	"singula.finance/node/chain"
)

func TestProcesser(t *testing.T) {

	ntchan := netio.ConnNtchanStub("test", "testout")
	netio.NetConnectorSetupMock(ntchan)

	m := netio.MessageJSON{MessageType: "REQ", Command: "PING"}
	jm, _ := json.Marshal(m)

	ntchan.Reader_queue <- string(jm)

	readout := <-ntchan.REQ_in

	if readout != string(jm) {
		t.Error("process error ", readout)
	}

	ntchan.Reader_queue <- string(jm)

	readout2 := <-ntchan.REQ_in

	mgr := chain.CreateManager()
	reply := RequestReply(mgr, ntchan, readout2)

	m = netio.MessageJSON{MessageType: "REP", Command: "PONG"}
	jm, _ = json.Marshal(m)

	if reply != string(jm) {
		t.Error("process reply error ", reply)
	}

}

func TestProcesserLoop(t *testing.T) {

	ntchan := netio.ConnNtchanStub("test", "testout")
	netio.NetConnectorSetupMock(ntchan)

	m := netio.MessageJSON{MessageType: "REQ", Command: "PING"}
	jm, _ := json.Marshal(m)

	ntchan.Reader_queue <- string(jm)

	readout2 := <-ntchan.REQ_in

	mgr := chain.CreateManager()
	reply := RequestReply(mgr, ntchan, readout2)

	m = netio.MessageJSON{MessageType: "REP", Command: "PONG"}
	jm, _ = json.Marshal(m)

	if reply != string(jm) {
		t.Error("process reply error ", reply)
	}

}

func TestBalanceReq(t *testing.T) {

	ntchan := netio.ConnNtchanStub("test", "testout")
	netio.NetConnectorSetupMock(ntchan)

	mgr := chain.CreateManager()
	mgr.InitAccounts()

	bb, _ := json.Marshal("P2e2bfb58c9db")
	raw := json.RawMessage(bb)

	m := netio.MessageJSON{MessageType: "REQ", Command: "BALANCE", Data: raw}
	jm, _ := json.Marshal(m)

	ntchan.Reader_queue <- string(jm)

	readout2 := <-ntchan.REQ_in

	reply := RequestReply(mgr, ntchan, readout2)

	cc, _ := json.Marshal(400)
	raw2 := json.RawMessage(cc)

	m = netio.MessageJSON{MessageType: "REP", Command: "BALANCE", Data: raw2}
	jm, _ = json.Marshal(m)

	if reply != string(jm) {
		t.Error("process reply error ", reply, raw2)
	}

}
func TestTxReq(t *testing.T) {

	ntchan := netio.ConnNtchanStub("test", "testout")
	netio.NetConnectorSetupMock(ntchan)

	mgr := chain.CreateManager()
	mgr.InitAccounts()

	keypair_sender := crypto.PairFromSecret("sender")
	pubkey_sender := crypto.PubKeyToHex(keypair_sender.PubKey)
	addr_sender := crypto.Address(pubkey_sender)

	keypair_receiver := crypto.PairFromSecret("receiver")
	addr_receiver := crypto.Address(crypto.PubKeyToHex(keypair_receiver.PubKey))

	// create tx
	tx := block.Tx{Nonce: 1, Amount: 10, Sender: addr_sender, Receiver: addr_receiver, SenderPubkey: pubkey_sender}
	// sign tx
	signature := block.SignTx(tx, keypair_sender.PrivKey)
	tx.Signature = hex.EncodeToString(signature.Serialize())

	bb, _ := json.Marshal(tx)
	raw := json.RawMessage(bb)
	fmt.Println(raw)

	if len(raw) == 0 {
		t.Error("raw ", raw)
	}

	//TODO parse back and check

	//TODO let mgr apply the tx

}
