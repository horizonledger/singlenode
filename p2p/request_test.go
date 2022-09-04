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
