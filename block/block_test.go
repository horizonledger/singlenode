package block

import (
	"testing"
	"time"
)

//basic block functions
func TestBasic(t *testing.T) {
	var txs []Tx
	var ph [32]byte
	new_block := Block{Height: 0, Txs: txs, Prev_Block_Hash: ph, Timestamp: time.Now()}
	if new_block.Height != 0 {
		t.Error("fail")
	}

	//TODO add tx to block

	//TODO check block calculations

	//TODO append block

	//TODO sign block

	//TODO check signature block
}
