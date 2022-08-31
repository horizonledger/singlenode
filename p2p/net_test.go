package p2p

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"singula.finance/netio"
)

const test_node_port = 8080

func initserver() *TCPNode {
	//log.Println("initserver")
	// Start the new server

	testsrv, err := NewNode()
	testsrv.addr = ":" + strconv.Itoa(test_node_port)

	if err != nil {
		log.Println("error starting TCP server")
		return testsrv
	} else {
		log.Println("start ", testsrv)
	}

	// Run the server in Goroutine to stop tests from blocking
	// test execution
	//log.Println("initserver  ", testsrv)

	go testsrv.RunTCP()
	//log.Println("waiting ", newpeerchan)
	go testsrv.HandleConnectTCP()

	return testsrv
}

func testclient() netio.Ntchan {
	time.Sleep(200 * time.Millisecond)
	addr := ":" + strconv.Itoa(test_node_port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		//t.Error("could not connect to server: ", err)
	}
	//t.Error("...")
	//log.Println("connected")
	ntchan := netio.ConnNtchan(conn, "client", addr, false)
	netio.NetConnectorSetupEcho(ntchan)
	defer conn.Close()
	return ntchan

}

func TestServer_Run(t *testing.T) {
	testsrv := initserver()
	defer testsrv.Close()

	time.Sleep(800 * time.Millisecond)

	// Simply check that the server is up and can accept connections
	ntclient := testclient()

	if ntclient.SrcName != "client" {
		t.Error("name")
	}
	// for ok := true; ok; testsrv.accepting = false {
	// 	log.Println(testsrv.accepting)
	// 	time.Sleep(100 * time.Millisecond)
	// }

	time.Sleep(1000 * time.Millisecond)
	//log.Println("TestServer_Run > ", testsrv, testsrv.Peers)

	if !testsrv.accepting {
		t.Error("not accepting")
	}

	if len(testsrv.Peers) != 1 {
		t.Error("no peers ", testsrv.Peers, len(testsrv.Peers))
	}

	ntclient2 := testclient()

	time.Sleep(1000 * time.Millisecond)

	if len(testsrv.Peers) != 2 {
		t.Error("no peers ", testsrv.Peers, len(testsrv.Peers))
	}

	if ntclient2.Alias == "xx" {
		t.Error("client nil")
	}
	// if ntclient2.Alias != "client" {
	// 	t.Error("name ", ntclient2.Alias)
	// }

	//ntclient.Writer_queue <- "test"

	//ntclient.quitchan <- true

	//TODO close time limit
	//https://stackoverflow.com/questions/49872097/idiomatic-way-for-reading-from-the-channel-for-a-certain-time
	//outread := <-testsrv.Peers[0].NTchan.Reader_queue

	// time.Sleep(1000 * time.Millisecond)

	// //ntclient2.Writer_queue <- "test"

	// testsrv.Close()

	// //outread := testsrv.peers.Reader_queue
	// // outread := <-testsrv.Peers[0].NTchan.Reader_queue

	// if outread != "test" {
	// 	t.Error("not read")
	// }

	// testsrv.Close()

}

// func TestServer_Write(t *testing.T) {

// 	testsrv := initserver()
// 	defer testsrv.Close()

// 	clientNt := testclient()
// 	go netio.ReadLoop(clientNt)

// 	time.Sleep(2000 * time.Millisecond)

// 	peers := testsrv.GetPeers()
// 	if len(peers) != 1 {
// 		t.Error("no peers ", testsrv.Peers, len(peers))
// 	}

// 	firstpeer := peers[0]

// 	if !xutils.IsEmpty(firstpeer.NTchan.Writer_queue, 1*time.Second) {
// 		t.Error("fail")
// 	}

// 	reqs := "hello world"
// 	n, err := netio.NetWrite(firstpeer.NTchan, reqs)

// 	if err != nil {
// 		t.Error("could not write to server:", err)
// 	}

// 	delimsize := 1
// 	l := len([]byte(reqs)) + delimsize
// 	if n != l {
// 		t.Error("wrong bytes written ", l)
// 	}

// 	time.Sleep(100 * time.Millisecond)

// 	rmsg1 := <-clientNt.Reader_queue
// 	if rmsg1 != "new peer connected. total peers 1" {
// 		t.Error("different message on reader ", rmsg1)
// 	}

// 	rmsg := <-clientNt.Reader_queue
// 	if rmsg != reqs {
// 		t.Error("different message on reader ", rmsg)
// 	}
// 	// if isEmpty(clientNt.Reader_queue, 1*time.Second) {
// 	// 	t.Error("fail")
// 	//}

// }

// func TestServer_WriteProcess(t *testing.T) {

// 	testsrv := initserver()
// 	defer testsrv.Close()

// 	clientNt := testclient()
// 	go netio.ReadLoop(clientNt)
// 	go netio.ReadProcessor(clientNt)

// 	time.Sleep(2000 * time.Millisecond)

// 	peers := testsrv.GetPeers()
// 	if len(peers) != 1 {
// 		t.Error("no peers ", testsrv.Peers, len(peers))
// 	}

// 	firstpeer := peers[0]

// 	if !xutils.IsEmpty(firstpeer.NTchan.Writer_queue, 1*time.Second) {
// 		t.Error("fail")
// 	}

// 	reqs := "REQ PING"
// 	n, err := netio.NetWrite(firstpeer.NTchan, reqs)

// 	if err != nil {
// 		t.Error("could not write to server:", err)
// 	}

// 	delimsize := 1
// 	l := len([]byte(reqs)) + delimsize
// 	if n != l {
// 		t.Error("wrong bytes written ", l)
// 	}

// 	time.Sleep(100 * time.Millisecond)

// 	rmsg1 := <-clientNt.Reader_queue
// 	if rmsg1 != "new peer connected. total peers 1" {
// 		t.Error("different message on reader ", rmsg1)
// 	}

// 	rmsg := <-clientNt.Reader_queue
// 	if rmsg != "pong" {
// 		t.Error("different message on reader ", rmsg)
// 	}
// 	// if isEmpty(clientNt.Reader_queue, 1*time.Second) {
// 	// 	t.Error("fail")
// 	//}

// }

func TestServer_Process(t *testing.T) {

	//bug: times out
	// testsrv := initserver()
	// defer testsrv.Close()

	// clientNt := testclient()
	// go netio.ReadLoop(clientNt)
	// go netio.ReadProcessorEcho(clientNt)
	// go netio.WriteProcessor(clientNt)
	// go netio.WriteLoop(clientNt, 300*time.Millisecond)

	// time.Sleep(2000 * time.Millisecond)

	// peers := testsrv.GetPeers()
	// if len(peers) != 1 {
	// 	t.Error("no peers ", testsrv.Peers, len(peers))
	// }

	// firstpeer := peers[0]

	// if !xutils.IsEmpty(firstpeer.NTchan.Writer_queue, 1*time.Second) {
	// 	t.Error("fail")
	// }

	// reqs := "hello world"
	// _, err := netio.NetWrite(firstpeer.NTchan, reqs)

	// if err != nil {
	// 	t.Error("could not write to server:", err)
	// }

	// time.Sleep(100 * time.Millisecond)

	// rmsg := <-clientNt.Reader_queue
	// if rmsg != reqs {
	// 	t.Error("different message on reader ", rmsg)
	// }

	// outmsg := <-clientNt.Writer_queue
	// if outmsg != "echo: "+reqs {
	// 	t.Error("different message on reader ", rmsg)
	// }
	// if isEmpty(clientNt.Reader_queue, 1*time.Second) {
	// 	t.Error("fail")
	//}

}
