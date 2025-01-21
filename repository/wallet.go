package repository

import (
	"context"
	"time"
	
	"blockchain-mining/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *Repository) CreateNewWallet(wallet *types.Wallet) error {
	ctx := context.Background()
	wallet.Time = uint64(time.Now().Unix())

	opt := options.Update().SetUpsert(true)

	filter := bson.M{"privateKey": wallet.PrivateKey}
	update := bson.M{"$set": wallet}

	if _, err := r.wallet.UpdateOne(ctx, filter, update, opt); err != nil {
		return err
	} else {
		return nil
	}
}
