package block

import (
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg/chainhash"

	"singula.finance/netio/crypto"
)

//"singula/node/block"

//potential TransactionTypes
// VOTE_DELEGATE
// REGISTER_NAME

const (
	CASH_TRANSFER = "CASH_TRANSFER"
	//DELEGATE_REGISTER = "DELEGATE_REGISTER"
	CONCESSION_REG = "CONCESSION_REG"
)

type TxSigmap struct {
	SenderPubkey string `json:"senderPubkey"`
	Signature    string `json:"signature"`
}

// type TxEnv struct {

// }

//TODO distinguish between tx and signed tx with extra struct
type Tx struct {
	TxType   string `json:"txType"`
	Amount   int    `json:"amount"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Nonce    int    `json:"nonce"`
	//TODO replace with txsig
	SenderPubkey string `json:"senderPubkey,omitempty"`
	Signature    string `json:"signature,omitempty"`
	//Id           [32]byte `edn:"id"`           //gets assigned when verified in a block

	//fee
	//txtype
	//timestamp

	//confirmations
	//height
}

// type SimpleTx struct {
// 	Amount   int    `edn:"amount"`
// 	Sender   string `edn:"sender"`   //[32]byte
// 	Receiver string `end:"receiver"` //[32]byte
// 	//Nonce        int    `edn:"Nonce"`
// }

// type TxSigmapEdn struct {
// 	SenderPubkey string `edn:"senderPubkey"`
// 	Signature    string `edn:"signature"`
// }

// // type TxExpr struct {
// // 	TxType   string   `edn:"TxType"`
// // 	Transfer SimpleTx `edn:"TxTransfer"`
// // 	Sigmap   TxSigmap `edn:"Sigmap"`
// // }

// type TxEdn struct {
// 	TxType   string `edn:"TxType"`
// 	Amount   int    `edn:"Amount"`
// 	Sender   string `edn:"Sender"`   //[32]byte
// 	Receiver string `end:"Receiver"` //[32]byte
// 	//TODO delete
// 	SenderPubkey string `edn:"SenderPubkey"` //hex string
// 	Signature    string `edn:"Signature"`    //hex string
// 	Nonce        int    `edn:"Nonce"`
// 	//Id           [32]byte `edn:"id"`           //gets assigned when verified in a block

// 	//fee
// 	//txtype
// 	//timestamp

// 	//confirmations
// 	//height
// }

func SignTx(tx Tx, privkey btcec.PrivateKey) btcec.Signature {
	//TODO sign tx not just id
	txJson, _ := json.Marshal(tx)
	//log.Println(string(txJson))
	//message := fmt.Sprintf("%d", tx.Id)

	messageHash := chainhash.DoubleHashB([]byte(txJson))
	signature, err := privkey.Sign(messageHash)
	if err != nil {
		fmt.Println("err ", err)
		//return
	}
	return *signature

}

func RemoveSigTx(tx Tx) Tx {
	tx.Signature = ""
	return tx
}

func RemovePubTx(tx Tx) Tx {
	tx.SenderPubkey = ""
	return tx
}

func VerifyTxSig(tx Tx) bool {
	pubkey := crypto.PubKeyFromHex(tx.SenderPubkey)
	sighex := tx.Signature
	sign := crypto.SignatureFromHex(sighex)
	//need to remove sig and pubkey for validation
	tx = RemoveSigTx(tx)
	tx = RemovePubTx(tx)

	txJson, _ := json.Marshal(tx)
	//log.Println("verify sig for tx ", string(txJson))
	verified := crypto.VerifyMessageSignPub(sign, pubkey, string(txJson))
	return verified

}

// // //hash of a transaction, currently sha256 of the nonce
// // //TODO hash properly
// // func TxHash(tx block.Tx) [32]byte {
// // 	b := []byte(string(tx.Nonce)[:])
// // 	hash := sha256.Sum256(b)
// // 	return hash
// // }

// // //sign tx and add signature and pubkey
// func SignTxAdd(tx block.Tx, keypair Keypair) block.Tx {

// 	signature := SignTx(tx, keypair.PrivKey)
// 	sighex := hex.EncodeToString(signature.Serialize())

// 	tx.Signature = sighex
// 	tx.SenderPubkey = PubKeyToHex(keypair.PubKey)
// 	return tx
// }
