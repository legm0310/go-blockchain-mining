package types

type Wallet struct {
	PrivateKey string `json:"privateKey" bson:"privateKey"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	Time       uint64 `json:"time" bson:"time"`
}

type Block struct {
	Time         int64          `json:"time" bson:"time"`
	Hash         []byte         `json:"hash" bson:"hash"`
	PrevHash     []byte         `json:"prevHash" bson:"prevHash"`
	Nonce        int64          `json:"nonce" bson:"nonce"`
	Height       int64          `json:"height" bson:"height"`
	Transactions []*Transaction `json:"transactions" bson:"transactions"`
}

type Transaction struct {
	Block   int64  `json:"block" bson:"block"`
	Time    int64  `json:"time" bson:"time"`
	From    string `json:"from" bson:"from"`
	To      string `json:"to" bson:"to"`
	Amount  string `json:"amount" bson:"amount"`
	Message string `json:"message" bson:"message"`
	Tx      string `json:"tx" bson:"tx"`
}
