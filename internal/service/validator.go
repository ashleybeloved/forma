package service

import "forma/internal/config"

type ValidatorService struct {
	Config *config.Config
}

func NewValidatorService(cfg *config.Config) *ValidatorService {
	return &ValidatorService{
		Config: cfg,
	}
}
