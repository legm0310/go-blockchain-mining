package service

import (
	"blockchain-mining/types"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (s *Service) CreateBlock(txs []*types.Transaction, prevHash []byte, height int64) *types.Block {
	var pHash []byte

	if latestBlock, err := s.repository.GetLatestBlock(); err == nil {
		if err == mongo.ErrNoDocuments {
			s.log.Info("Genesis Block Will Be Created")
			//Create Genesis Block
			newBlock := createBlockInner(txs, pHash, height)
			pow := s.NewPow(newBlock)
			newBlock.Nonce, newBlock.Hash = pow.RunMining()
			return newBlock
			// create new block
			// mining
		} else {
			s.log.Crit("Failed to get latest block", "err", err)
			panic(err)
		}
	} else {
		pHash = latestBlock.Hash
		newBlock := createBlockInner(txs, pHash, height)
		pow := s.NewPow(newBlock)
		newBlock.Nonce, newBlock.Hash = pow.RunMining()
		// create new block
		// mining
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
