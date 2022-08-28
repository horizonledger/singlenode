package node

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
)

//TODO move to node implementation

func RequestReplyFun(ntchan netio.Ntchan, msg netio.MessageJSON) netio.MessageJSON {

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
		//TODO
		//balance := t.Mgr.State.Accounts[a]
		//fmt.Println("balance for ", a, balance, t.Mgr.State.Accounts)

		balance := 100
		balJson, _ := json.Marshal(balance)
		raw := json.RawMessage(balJson)

		//rmsg := MessageJSON{MessageType: "REP", Command: "BALANCE", Data: &balJson}
		rmsg := netio.MessageJSON{MessageType: "REP", Command: "BALANCE", Data: raw}

		return rmsg

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

func RequestReply(ntchan netio.Ntchan, msgString string) string {

	//TODO separate namespace

	msg, _ := netio.ParseLineJson(msgString)

	rmsg := RequestReplyFun(ntchan, msg)
	rmsgstr, _ := json.Marshal(rmsg)

	fmt.Sprintf("Handle cmd %v", msg.Command)

	return string(rmsgstr)
}
