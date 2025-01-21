package repository

import (
	"context"
	"time"

	"blockchain-mining/config"

	"github.com/inconshreveable/log15"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client

	wallet *mongo.Collection
	tx     *mongo.Collection

	config *config.Config
	log    log15.Logger
}

func NewRepository(config *config.Config) (*Repository, error) {
	r := &Repository{
		config: config,
		log:    log15.New("module", "repository"),
	}

	var err error
	ctx := context.Background()

	if r.client, err = mongo.Connect(ctx, options.Client().ApplyURI(config.Mongo.Uri)); err != nil {
		r.log.Error("Failed to connect to mongo", "config", config.Mongo.Uri, "err", err)
		return nil, err
	} else if err = r.client.Ping(ctx, nil); err != nil {
		r.log.Error("Failed to ping mongo", "err", err)
		return nil, err
	} else {
		db := r.client.Database(config.Mongo.DB, nil)

		r.wallet = db.Collection("wallet")
		r.tx = db.Collection("tx")

		r.log.Info("Success To Connect Repository", "info", time.Now().Unix(), "repository", config.Mongo.Uri, "db", db)
		return r, nil
	}
}
