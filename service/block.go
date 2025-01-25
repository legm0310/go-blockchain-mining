package service

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hacpy/go-ethereum/common"
	"github.com/hacpy/go-ethereum/crypto"

	"blockchain-mining/types"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(from, to, value string) {
	var block *types.Block
	if latestBlock, err := s.repository.GetLatestBlock(); err == nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will Be Created")
			genesisMessage := "This is First Genesis Block"
			// TODO-> privateKey

			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx := createTransaction(genesisMessage, common.Address{}.String(), pk, to, value, 1)
				block = createBlockInner([]*types.Transaction{tx}, "", 1)
			}
		}
	} else {
		var tx *types.Transaction
		if common.HexToAddress(from) == (common.Address{}) {
			// Mint
			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx = createTransaction("MintCoin", common.Address{}.String(), pk, to, value, 1)
			}
		} else {
			// Transfer
			if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
				panic(err)
			} else {
				// TODO -> from의 밸런스 체크
				tx = createTransaction("TransferCoin", from, wallet.PrivateKey, to, value, 1)
			}
		}

		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)
	}

	pow := s.NewPow(block)
	block.Nonce, block.Hash = pow.RunMining()

	if err := s.repository.SaveBlock(block); err != nil {
		panic(err)
	}
}

func createBlockInner(txs []*types.Transaction, prevHash string, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         "",
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, pk, to, amount string, block int64) *types.Transaction {
	data := struct {
		Message string `json:"message"`
		From    string `json:"from"`
		To      string `json:"to"`
		Amount  string `json:"amount"`
	}{
		Message: message,
		From:    from,
		To:      to,
		Amount:  amount,
	}
	fmt.Println(pk)
	dataToSign := fmt.Sprintf("%x\n", data)

	if ecdsaPrivateKey, err := crypto.HexToECDSA(pk[2:]); err != nil {
		panic(err)
	} else if r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivateKey, []byte(dataToSign)); err != nil {
		panic(err)
	} else {
		signature := append(r.Bytes(), s.Bytes()...)

		return &types.Transaction{
			Block:   block,
			Time:    time.Now().Unix(),
			From:    from,
			To:      to,
			Amount:  amount,
			Message: message,
			Tx:      hex.EncodeToString(signature),
		}
	}
}

func HashTransactions(b *types.Block) []byte {
	var txHashes [][]byte

	for _, tx := range b.Transactions {
		var encoded bytes.Buffer

		enc := gob.NewEncoder(&encoded)

		if err := enc.Encode(tx); err != nil {
			panic(err)
		} else {
			txHashes = append(txHashes, encoded.Bytes())
		}
	}
	tree := NewMerkleTree(txHashes)
	return tree.RootNode.Data
}
