package p2p

import (
	"encoding/json"
	"testing"

	"singula.finance/netio"
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
