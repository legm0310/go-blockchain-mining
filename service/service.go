package service

import (
	"blockchain-mining/config"
	"blockchain-mining/repository"
)

type Service struct {
	config *config.Config

	repository *repository.Repository
}

func NewService(config *config.Config, repository *repository.Repository) *Service {
	s := &Service{
		config:     config,
		repository: repository,
	}

	return s
}
