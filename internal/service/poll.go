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
	GeoIP  *GeoIPService
}

func NewPollService(repo *repository.PollRepository, cfg *config.Config, geoIP *GeoIPService) *PollService {
	return &PollService{
		Repo:   repo,
		Config: cfg,
		GeoIP:  geoIP,
	}
}

func (s *PollService) GetPollByShortID(shortID string) (*model.Poll, error) {
	poll, err := s.Repo.GetPollByShortID(shortID)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) GetAllMyPolls(creatorID, limit, offset int) ([]model.Poll, error) {
	polls, err := s.Repo.GetPollsByCreatorID(creatorID, limit, offset)
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

func (s *PollService) UpdatePoll(id int, title, description string, config model.PollConfig, userID int) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return ErrMarshalJSON
	}

	poll := &model.Poll{
		ID:          id,
		Title:       title,
		Description: description,
	}

	return s.Repo.UpdatePoll(poll, configBytes, userID)
}

func (s *PollService) DeletePoll(id, creatorID int) error {
	return s.Repo.DeletePoll(id, creatorID)
}

func (s *PollService) Vote(tokenStr, pollShortID, guestToken, ip string, answers []model.Answer) error {
	var userID int

	if tokenStr == "" {
		userID = -1
	} else {
		claims, err := pkg.ValidateToken(tokenStr, s.Config.JWTSecretKey)
		if err != nil {
			return ErrInvalidToken
		}

		userID = claims.UserID
	}

	countryCode := s.GeoIP.GetCountryCodeFromIP(ip)

	vote := &model.Vote{
		PollShortID: pollShortID,
		UserID:      userID,
		IP:          ip,
		CountryCode: countryCode,
		GuestToken:  guestToken,
	}

	// Mock "secured" field for a time
	voted, err := s.Repo.HasVoted(false, vote.PollShortID, vote.IP, vote.GuestToken, vote.UserID)
	if err != nil {
		return err
	}

	if voted {
		return ErrAlreadyVoted
	}

	err = s.Repo.Vote(vote, answers)
	if err != nil {
		return err
	}

	return nil
}

func (s *PollService) CheckVote(tokenStr, pollShortID, guestToken, ip string) (bool, error) {
	var userID int

	if tokenStr == "" {
		userID = -1
	} else {
		claims, err := pkg.ValidateToken(tokenStr, s.Config.JWTSecretKey)
		if err != nil {
			return true, ErrInvalidToken
		}

		userID = claims.UserID
	}

	// Mock "secured" field for a time
	return s.Repo.HasVoted(false, pollShortID, ip, guestToken, userID)
}
