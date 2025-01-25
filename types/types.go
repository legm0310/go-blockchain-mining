package types

type Wallet struct {
	PrivateKey string `json:"privateKey" bson:"privateKey"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	Time       uint64 `json:"time" bson:"time"`
}

// 일반적으로 height는 한 블록의 트랜잭션 최대 갯수, 여기선 블록 순서
type Block struct {
	Time         int64          `json:"time" bson:"time"`
	Hash         string         `json:"hash" bson:"hash"`
	PrevHash     string         `json:"prevHash" bson:"prevHash"`
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
