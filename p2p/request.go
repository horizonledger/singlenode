package p2p

// network communication layer (netio)

// netio -> semantics of channels
// TCP/IP -> golang net

// TODO
// create a channel wrapper struct
// which has a flag if its in or out flow
// see whitepaper for details
// type Nchain {
// c chan string
// inflow }

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"time"

	"singula.finance/netio"
	"singula.finance/node/chain"
)

//TODO pass manager

func RequestReplyFun(mgr chain.ChainManager, ntchan netio.Ntchan, msg netio.MessageJSON) netio.MessageJSON {

	switch msg.Command {

	case netio.CMD_PING:
		rmsg := netio.MessageJSON{MessageType: "REP", Command: "PONG"}
		return rmsg

	case netio.CMD_EXIT:
		//TODO close chans?
		err := ntchan.Conn.(*net.TCPConn).SetLinger(0)
		if err != nil {
			log.Printf("Error when setting linger: %s", err)
		} else {
			fmt.Println("connection closed")
			//quite all
			//TODO fix
			//ntchan.quitchan <- true

		}

	case netio.CMD_TIME:
		dt := time.Now()
		dtlJson, _ := json.Marshal(dt.String())
		r := json.RawMessage(dtlJson)

		rmsg := netio.MessageJSON{MessageType: "REP", Command: "TIME", Data: r}
		return rmsg

	case netio.CMD_BALANCE:
		var a string
		json.Unmarshal(msg.Data, &a)
		b := string(a)
		balance := mgr.State.Accounts[b]

		balJson, _ := json.Marshal(balance)
		raw := json.RawMessage(balJson)

		rmsg := netio.MessageJSON{MessageType: "REP", Command: "BALANCE", Data: raw}

		return rmsg

	case netio.CMD_TX:
		fmt.Println("Handle tx")
	// msg = HandleTx(t, msg)
	// data, _ := json.Marshal(msg.Data)
	// reply_msg = netio.EdnConstructMsgMapData(netio.REP, netio.CMD_GETBLOCKS, string(data))

	// case CMD_REGISTERALIAS:
	// 	//TODO only pointer is set
	// 	ntchan.Alias = "123"
	// 	reply_msg := fmt.Sprintf("new alias %v", ntchan.Alias)
	// 	//fmt.Printf("new alias %v", ntchan.Alias)
	// 	return reply_msg

	//handshake
	// case CMD_REGISTERPEER:
	// 	reply_msg := "todo"
	// 	return reply_msg

	default:
		errormsg := "Error: not found command"
		fmt.Println(errormsg)
		//xjson, _ := json.Marshal("")
		msg := netio.MessageJSON{MessageType: netio.REP, Command: netio.CMD_ERROR}
		return msg
		//reply_msg = ToJSONMessage(msg)
	}

	return netio.MessageJSON{MessageType: netio.REP, Command: netio.CMD_ERROR}

}

func RequestReply(mgr chain.ChainManager, ntchan netio.Ntchan, msgString string) string {

	msg, err := netio.ParseLineJson(msgString)

	// fmt.Println("msg >> ", msg, err)

	if err == nil {
		fmt.Sprintf("Handle cmd %v", msg.Command)
		rmsg := RequestReplyFun(mgr, ntchan, msg)
		fmt.Sprintf("return msg %v", rmsg)
		rmsgstr, _ := json.Marshal(rmsg)

		return string(rmsgstr)
	} else {
		fmt.Println("?????")
		returnerr := fmt.Sprintf("error %v", err)
		fmt.Println(returnerr)
		return returnerr
	}
}

func RequestReplyLoop(mgr chain.ChainManager, ntchan netio.Ntchan) {
	for {
		msg := <-ntchan.REQ_in
		//fmt.Println("request %v", msg)
		//vlog(ntchan, "request "+msg)
		reply := RequestReply(mgr, ntchan, msg)
		//reply := "testing"
		ntchan.REP_out <- reply
	}
}
