package chain

import (
	"testing"

	"golang.org/x/exp/maps"
)

func TestBalance(t *testing.T) {

	mgr := CreateManager()
	mgr.InitAccounts()

	ks := maps.Keys(mgr.State.Accounts)

	if mgr.State.Accounts[ks[0]] == 0 {
		t.Error("0 balance")
	}

	genBlock := MakeGenesisBlock()
	if genBlock.Height != 0 {
		t.Error("wrong height")
	}

	mgr.ApplyBlock(genBlock)

	// if reply.MessageType != netio.REP {
	// 	t.Error("balance msg", reply_msg)
	// }

}
