package main

import (
	"encoding/json"
	"testing"

	"singula.finance/netio"
)

func TestProcesser(t *testing.T) {

	ntchan := netio.ConnNtchanStub("test", "testout")
	netio.NetConnectorSetupMock(ntchan)

	m := netio.MessageJSON{MessageType: "REQ", Command: "PING"}
	jm, _ := json.Marshal(m)
	if string(jm) !=
		"{\"messagetype\":\"REQ\",\"command\":\"PING\"}" {
		t.Error("encoding error: ", string(jm))
	}

	go func() { ntchan.Reader_queue <- string(jm) }()

	readout := <-ntchan.REQ_in

	if readout != "{\"messagetype\":\"REQ\",\"command\":\"PING\"}" {
		t.Error("process error ", readout)
	}

	go func() { ntchan.Reader_queue <- string(jm) }()

	readout2 := <-ntchan.REQ_in

	if readout2 == "{\"messagetype\":\"REQ\",\"command\":\"PING\"}" {
		//ntchan.REP_out <- "REP PONG"
	}

	reply := RequestReply(ntchan, readout2)

	if reply != "{\"messagetype\":\"REP\",\"command\":\"PONG\"}" {
		t.Error("process reply error ", reply)
	}

}
