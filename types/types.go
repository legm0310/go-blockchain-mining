package types

type Wallet struct {
	PrivateKey string `json:"privateKey" bson:"privateKey"`
	PublicKey  string `json:"publicKey" bson:"publicKey"`
	Time       uint64 `json:"time" bson:"time"`
}

type Tx struct {
	Hash      string `json:"hash" bson:"hash"`
	PrevHash  string `json:"prevHash" bson:"prevHash"`
	Timestamp uint64 `json:"timestamp" bson:"timestamp"`
}
