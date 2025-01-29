package service

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"blockchain-mining/types"
	"github.com/hacpy/go-ethereum/common"
	"github.com/hacpy/go-ethereum/crypto"
	"github.com/shopspring/decimal"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateBlock(from, to, value string) {
	var block *types.Block
	toBalance := "0"

	if latestBlock, err := s.repository.GetLatestBlock(); err != nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will Be Created")
			genesisMessage := "This is First Genesis Block"

			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx := createTransaction(genesisMessage, common.Address{}.String(), pk, to, value, 1)
				block = createBlockInner([]*types.Transaction{tx}, "", 1)
			}
		} else {
			panic(err)
		}
	} else {
		var tx *types.Transaction

		if common.HexToAddress(from) == (common.Address{}) {
			// Mint
			if pk, _, err := s.newKeyPair(); err != nil {
				panic(err)
			} else {
				tx = createTransaction("MintCoin", common.Address{}.String(), pk, to, value, 1)
				toBalance = value
			}
		} else {
			// Transfer
			if wallet, err := s.repository.GetWalletByPublicKey(from); err != nil {
				panic(err)
			} else if toWallet, err := s.repository.GetWalletByPublicKey(to); err != nil {
				if err == mongo.ErrNoDocuments {
					s.log.Debug("Can't Find To Wallet", "to", to)
				} else {
					panic(err)
				}
				return
			} else {
				//from 주소에서 Balance가 있는지
				fromDecimalBalance, _ := decimal.NewFromString(wallet.Balance)
				valueDecimal, _ := decimal.NewFromString(value)
				toDecimalBalance, _ := decimal.NewFromString(toWallet.Balance)

				if fromDecimalBalance.Cmp(valueDecimal) == -1 {
					s.log.Debug("Failed to transfer coin by From Balance", "from", from, "balance", wallet.Balance, "value", value)
					return
				} else {
					toDecimalBalance = toDecimalBalance.Add(valueDecimal)
					toBalance = toDecimalBalance.String()

					fromDecimalBalance = fromDecimalBalance.Sub(valueDecimal)
					value = fromDecimalBalance.String()
				}

				tx = createTransaction("TransferCoin", from, wallet.PrivateKey, to, value, 1)
			}
		}
		block = createBlockInner([]*types.Transaction{tx}, latestBlock.Hash, latestBlock.Height+1)
	}
	pow := s.NewPow(block)
	block.Nonce, block.Hash = pow.RunMining()

	if err := s.repository.UpsertWalletsWhenTransfer(from, to, value, toBalance); err != nil {
		panic(err)
	} else if err := s.repository.SaveBlock(block); err != nil {
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

	dataToSign := fmt.Sprintf("%x\n", data)

	if ecdsaPrivateKey, err := crypto.HexToECDSA(strings.TrimPrefix(pk, "0x")); err != nil {
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
