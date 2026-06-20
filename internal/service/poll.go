package service

import (
	"encoding/json"
	"forma/internal/config"
	"forma/internal/model"
	"forma/internal/pkg"
	"forma/internal/repository"
	"time"
)

type PollService struct {
	Repo   *repository.PollRepository
	Config *config.Config
}

func NewPollService(repo *repository.PollRepository, cfg *config.Config) *PollService {
	return &PollService{
		Repo:   repo,
		Config: cfg,
	}
}

func (s *PollService) GetPollByShortID(shortID string) (*model.Poll, error) {
	poll, err := s.Repo.GetPollByShortID(shortID)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) GetAllMyPolls(creatorID int) ([]model.Poll, error) {
	polls, err := s.Repo.GetPollsByCreatorID(creatorID)
	if err != nil {
		return nil, err
	}

	return polls, nil
}

func (s *PollService) CreatePoll(title, description string, config model.PollConfig, creatorID int) (*model.Poll, error) {
	shortID := pkg.GenerateShortID(s.Config.ShortIDLength)

	configBytes, err := json.Marshal(config)
	if err != nil {
		return nil, ErrMarshalJSON
	}

	poll := &model.Poll{
		Title:       title,
		Description: description,
		Config:      config,
		ShortID:     shortID,
		CreatorID:   creatorID,
		CreatedAt:   time.Now(),
	}

	err = s.Repo.CreatePoll(poll, configBytes)
	if err != nil {
		return nil, err
	}

	return poll, nil
}
