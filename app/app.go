package app

import (
	"blockchain-mining/config"
	"blockchain-mining/repository"
	"blockchain-mining/service"
)

type App struct {
	config *config.Config

	service    *service.Service
	repository *repository.Repository
}

func NewApp(config *config.Config) {
	a := &App{
		config: config,
	}

	var err error

	if a.repository, err = repository.NewRepository(a.config); err != nil {
		panic(err)
	} else {
		a.service = service.NewService(a.config, a.repository)
	}
}
