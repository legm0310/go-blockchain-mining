package service

import (
	"blockchain-mining/config"
	"blockchain-mining/repository"
	
	"github.com/inconshreveable/log15"
)

type Service struct {
	config *config.Config

	repository *repository.Repository

	log log15.Logger
}

func NewService(config *config.Config, repository *repository.Repository) *Service {
	s := &Service{
		config:     config,
		repository: repository,
		log:        log15.New("module", "service"),
	}

	return s
}
