package service

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/hacpy/go-ethereum/crypto"

	"blockchain-mining/types"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte

	if latestBlock, err := s.repository.GetLatestBlock(); err == nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will Be Created")

			genesisMessage := "This is First Genesis Block"
			tx := createTransaction(genesisMessage, "0x7ce062a2aae399c6db5014ef682af37ff1ba0d58", "", "", 1)

			newBlock := createBlockInner([]*types.Transaction{tx}, pHash, height)
			pow := s.NewPow(newBlock)

			newBlock.Nonce, newBlock.Hash = pow.RunMining()

			return newBlock
		} else {
			s.log.Crit("Failed to get latest block", "err", err)
			panic(err)
		}
	} else {
		pHash = latestBlock.Hash

		newBlock := createBlockInner(txs, pHash, height)
		pow := s.NewPow(newBlock)

		newBlock.Nonce, newBlock.Hash = pow.RunMining()
		
		return newBlock
	}
}

func createBlockInner(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	return &types.Block{
		Time:         time.Now().Unix(),
		Hash:         []byte{},
		Transactions: txs,
		PrevHash:     prevHash,
		Nonce:        0,
		Height:       height,
	}
}

func createTransaction(message, from, to, amount string, block int64) *types.Transaction {
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

	dataToSign := fmt.Sprintf("%x\n", data)

	pk := "b18c7a77961e7dde666b6caa8aedae1a8f7e802ea985fbb9a49ee4d0fc7bdb86"
	if ecdsaPrivateKey, err := crypto.HexToECDSA(pk); err != nil {
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
