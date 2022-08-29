package chain

import (
	"testing"
)

func TestBalance(t *testing.T) {

	mgr := CreateManager()
	//node.Mgr = &mgr
	mgr.InitAccounts()
	//initSyncChain(node.Config)

	//genBlock := chain.MakeGenesisBlock()
	//mgr.ApplyBlock(genBlock)

	// if reply.MessageType != netio.REP {
	// 	t.Error("balance msg", reply_msg)
	// }

}
