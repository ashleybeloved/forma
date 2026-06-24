package service

import (
	"encoding/json"
	"forma/internal/config"
	"forma/internal/model"
	"forma/internal/pkg"
	"forma/internal/repository"
)

type PollService struct {
	Repo      *repository.PollRepository
	Config    *config.Config
	GeoIP     *GeoIPService
	Validator *ValidatorService
}

func NewPollService(repo *repository.PollRepository, cfg *config.Config, geoIP *GeoIPService, validator *ValidatorService) *PollService {
	return &PollService{
		Repo:      repo,
		Config:    cfg,
		GeoIP:     geoIP,
		Validator: validator,
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

func (s *PollService) CreatePoll(title, description string, config model.PollConfig, userID int, secured, authOnly bool) (*model.Poll, error) {
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
		CreatorID:   userID,
		Secured:     secured,
		AuthOnly:    authOnly,
	}

	err = s.Repo.CreatePoll(poll, configBytes)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func (s *PollService) UpdatePoll(short_id, title, description string, config model.PollConfig, userID int, secured, authOnly bool) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return ErrMarshalJSON
	}

	poll := &model.Poll{
		ShortID:     short_id,
		Title:       title,
		Description: description,
		Secured:     secured,
		AuthOnly:    authOnly,
	}

	return s.Repo.UpdatePoll(poll, configBytes, userID)
}

func (s *PollService) DeletePoll(short_id string, creatorID int) error {
	return s.Repo.DeletePoll(short_id, creatorID)
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

	poll, err := s.Repo.GetPollByShortID(pollShortID)
	if err != nil {
		return err
	}

	err = s.Validator.ValidateAnswers(poll.Config.Questions, answers)
	if err != nil {
		return err
	}

	if poll.AuthOnly && userID == -1 {
		return ErrAuthOnly
	}

	countryCode := s.GeoIP.GetCountryCodeFromIP(ip)

	vote := &model.Vote{
		PollShortID: pollShortID,
		UserID:      userID,
		IP:          ip,
		CountryCode: countryCode,
		GuestToken:  guestToken,
	}

	voted, err := s.Repo.HasVoted(poll.Secured, vote.PollShortID, vote.IP, vote.GuestToken, vote.UserID)
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

	poll, err := s.Repo.GetPollByShortID(pollShortID)
	if err != nil {
		return true, err
	}

	if poll.AuthOnly && userID == -1 {
		return true, ErrAuthOnly
	}

	return s.Repo.HasVoted(poll.Secured, pollShortID, ip, guestToken, userID)
}

func (s *PollService) GetPollStats(userID int, pollShortID string) (*model.Stats, error) {
	creator := s.Repo.IsCreator(userID, pollShortID)
	if !creator {
		return nil, ErrNotUserPoll
	}

	stats, err := s.Repo.GetPollStats(pollShortID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
