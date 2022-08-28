package node

//node.go is the main software which validators run

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"singula/node/chain"

	"golang.org/x/exp/maps"
	"singula.finance/netio"
)

// var blocktime = 10000 * time.Millisecond
// var logfile_name = "node.log"

// const LOGLEVEL_OFF = 0
const LOGLEVEL_ON = 1

type Config struct {
	NodeAlias     string
	Verbose       bool
	NodePort      int
	CreateGenesis bool
}

// //TODO rename
type TCPNode struct {
	NodePort      int
	Name          string
	addr          string
	server        net.Listener
	accepting     bool
	ConnectedChan chan net.Conn //channel of newly connected clients/peers
	Peers         []netio.Peer
	Mgr           *chain.ChainManager
	BROAD_out     chan string
	BROAD_in      chan string
	BROAD_signal  chan string
	Starttime     time.Time
	Logger        *log.Logger
	Loglevel      int
	Config        Config
	//
	ChatSubscribers []netio.Ntchan
}

// "DelegateName": "localhost",
// "PeerAddresses": ["polygonnode.com", "polygon.cc", "swix.io"],
// DelgateEnabled true
// CreateGenesis true
// Verbose true

// func (t *TCPNode) GetPeers() []netio.Peer {
// 	if &t.Peers == nil {
// 		return nil
// 	}
// 	return t.Peers
// }

// func (t *TCPNode) log(s string) {
// 	//fmt.Println(t.Loglevel)
// 	if t.Loglevel > LOGLEVEL_OFF {
// 		//t.Logger.Println(s)
// 		fmt.Println(s)
// 	}
// }

// start listening on tcp and handle connection through channels
func (t *TCPNode) RunTCP() (err error) {
	t.Starttime = time.Now()

	log.Println("node listens on " + t.addr)
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		//return errors.Wrapf(err, "Unable to listen on port %s\n", t.addr)
	}
	//run forever and don't close
	//defer t.Close()

	for {
		t.accepting = true
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			break
		}
		if conn == nil {
			err = errors.New("could not create connection")
			break
		}

		//t.log(fmt.Sprintf("new conn accepted %v", conn))
		//we put the new connection on the chan and handle there
		t.ConnectedChan <- conn

		// 	//TODO check if peers are alive see
		// 	//https://stackoverflow.com/questions/12741386/how-to-know-tcp-connection-is-closed-in-net-package
		// 	//https://gist.github.com/elico/3eecebd87d4bc714c94066a1783d4c9c

	}
	//t.log("end run")
	return
}

// func (t *TCPNode) HandleDisconnect() {

// }

//handle new connection
func (t *TCPNode) HandleConnectTCP() {

	//TODO! hearbeart, check if peers are alive
	//TODO! handshake
	log.Println("HandleConnectTCP")

	for {
		newpeerConn := <-t.ConnectedChan
		strRemoteAddr := newpeerConn.RemoteAddr().String()
		//t.log(fmt.Sprintf("accepted conn %v %v", strRemoteAddr, t.accepting))
		//t.log(fmt.Sprintf("new peer %v ", newpeerConn))
		// log.Println("> ", t.Peers)
		// log.Println("# peers ", len(t.Peers))
		//t.log(fmt.Sprintf("setup channels"))
		Verbose := true
		srcName := "localNode"
		destName := strRemoteAddr
		//ntchan := netio.ConnNtchan(newpeerConn, srcName, destName, Verbose, t.BROAD_signal)
		ntchan := netio.ConnNtchan(newpeerConn, srcName, destName, Verbose)

		rand.Seed(time.Now().UnixNano())
		ran := rand.Intn(100)
		ranname := fmt.Sprintf("ranPeer%v", ran)
		p := netio.Peer{Address: strRemoteAddr, NodePort: t.NodePort, NTchan: ntchan, Name: ranname}
		//t.log(fmt.Sprintf("new peer %v : %v", p.Name, p))
		t.Peers = append(t.Peers, p)

		go t.handleConnection(p)

		//conn.Close()

	}
}

//init an output connection
//TODO add this function
// check if connected inbound already
// func initOutbound(mainPeerAddress string, node_port int, verbose bool, BROAD_signal chan string) netio.Ntchan {
// 	log.Println("initOutbound")

// 	addr := mainPeerAddress + ":" + strconv.Itoa(node_port)
// 	//log.Println("dial ", addr)
// 	conn, err := net.Dial("tcp", addr)
// 	if err != nil {
// 		//log.Println("cant run")
// 		//return
// 	}

// 	//log.Println("connected")NetMsgRead
// 	//ntchan := netio.ConnNtchan(conn, "client", addr, verbose, BROAD_signal)
// 	ntchan := netio.ConnNtchan(conn, "client", addr, verbose)

// 	go netio.ReadLoop(ntchan)
// 	go netio.ReadProcessor(ntchan)
// 	go netio.WriteProcessor(ntchan)
// 	go netio.WriteLoop(ntchan, 300*time.Millisecond)
// 	return ntchan

// }

// //func (t *TCPNode) handleConnection(mgr *chain.ChainManager, ntchan netio.Ntchan) {
//func (t *TCPNode) handleConnection(mgr *chain.ChainManager, peer netio.Peer) {
func (t *TCPNode) handleConnection(peer netio.Peer) {
	//tr := 100 * time.Millisecond
	//defer ntchan.Conn.Close()
	//t.log(fmt.Sprintf("handleConnection"))

	//t.log(fmt.Sprintf("number of peers %v", len(t.Peers)))

	netio.NetConnectorSetup(peer.NTchan, RequestReply)

	//TODO
	//EXAMPLE
	//publoop all
	for i, peer := range t.Peers {
		fmt.Println(i)
		//fmt.Println("peer ", reflect.TypeOf(peer))
		c, _ := json.Marshal(len(t.Peers))
		msg := netio.Message{MessageType: "PUB", Command: netio.CMD_NUMPEERS, Data: c}

		peer.NTchan.Writer_queue <- netio.ToJSONMessage(msg)
	}

	//add to list of peers?

}

// // create a new node
func NewNode() (*TCPNode, error) {
	return &TCPNode{
		//addr:          addr,
		accepting:     false,
		ConnectedChan: make(chan net.Conn),
		Loglevel:      LOGLEVEL_ON,
	}, nil
}

// Close shuts down the TCP Server
func (t *TCPNode) Close() (err error) {
	return t.server.Close()
}

func runNode(t *TCPNode) {
	fmt.Printf("run node")

	go t.HandleConnectTCP()
	go t.RunTCP()

	// 	if err != nil {
	// 		node.log(fmt.Sprintf("error creating TCP server"))
	// 		return
	// 	}

	// 	//TODO! not get full chain after the init sync
	// 	// if !config.CreateGenesis {
	// 	// 	go func() {
	// 	// 		for {
	// 	// 			log.Println("fetch blocks loop")
	// 	// 			FetchAllBlocks(config, node)
	// 	// 			time.Sleep(10000 * time.Millisecond)
	// 	// 		}
	// 	// 	}()
	// 	// }

}

//WIP currently in testnet there is a single initiator which is the delegate expected to create first block
//TODO! replace with quering for blockheight?
func (t *TCPNode) initSyncChain(config Config) {
	if config.CreateGenesis {
		fmt.Println("CreateGenesis")
		genBlock := chain.MakeGenesisBlock()
		t.Mgr.ApplyBlock(genBlock)
		//TODO!
		t.Mgr.AppendBlock(genBlock)
		fmt.Println("accounts\n ", t.Mgr.State.Accounts)
		//keys := reflect.ValueOf(t.Mgr.State.Accounts).MapKeys()

		ks := maps.Keys(t.Mgr.State.Accounts)
		fmt.Println("all balances: ", t.Mgr.State.Accounts)
		fmt.Println("balance: ", t.Mgr.State.Accounts[ks[0]])

	}

	// else {

	// 	//TODO! apply blocks
	// 	success := t.Mgr.ReadChain()
	// 	t.log(fmt.Sprintf("read chain success %v", success))
	// 	loaded_height := len(t.Mgr.Blocks)
	// 	t.log(fmt.Sprintf("block height %d", loaded_height))

	// 	//TODO! age of latest block compared to local time
	// 	are_behind := loaded_height < 2
	// 	if are_behind {
	// 		t.Mgr.ResetBlocks()
	// 		log.Println("blocks after reset ", len(t.Mgr.Blocks))
	// 		FetchAllBlocks(config, t)
	// 	}

	// }
}

func getConfig() *Config {

	conffile := "config.json"
	log.Println("config file ", conffile)

	if _, err := os.Stat(conffile); os.IsNotExist(err) {
		log.Println("config file does not exist. create a file named ", conffile)
		return nil
	}

	content, err := ioutil.ReadFile(conffile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return &config

}

func runNodeWithConfig() {

	// 	go runAll(config)

	log.Println("main")

	config := getConfig()

	// 	node.setupLogfile()
	// 	node.log(fmt.Sprintf("PeerAddresses: %v", node.Config.PeerAddresses))

	log.Println("runNodeAll with config ", config)
	log.Println("verbose ", config.Verbose)

	node, err := NewNode()
	node.Config = *config

	node.addr = ":" + strconv.Itoa(config.NodePort)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	if err == nil {
		fmt.Printf("...")
		mgr := chain.CreateManager()
		node.Mgr = &mgr
		node.Mgr.InitAccounts()
		node.initSyncChain(node.Config)
		runNode(node)

		//TODO signatures of genesis
		fmt.Printf("InitAccounts")

	} else {
		log.Printf("error %v", err)
	}

	<-quit // This will block until you manually exists with CRl-C
	log.Println("\nnode exiting")
	// 	// log.Println("Got quit signal: shutdown node ...")
	//fmt.Printf("\nuser exit")
	// 	// signal.Reset(os.Interrupt)

	// 	//handle shutdown should never happen, need restart on OS level and error handling

}

func main() {

	m2 := netio.ConstructMsgSimple("REQ", "TIME")
	log.Println(m2)

	runNodeWithConfig()

	//GitCommit := os.Getenv("GIT_COMMIT")
	//fmt.Printf("--- run horizon ---\ngit commit: %s ----\n", GitCommit)

	//runNodeWithConfig()

}
