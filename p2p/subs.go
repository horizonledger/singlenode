package p2p

import (
	"fmt"

	"singula.finance/netio"
	"singula.finance/node/chain"
)

func PubSubLoop(mgr *chain.ChainManager, ntchan netio.Ntchan) {
	for {
		msg := <-ntchan.PUB_in
		fmt.Println("PUB %v", msg)

		//ntchan.REP_out <- reply
	}
}
