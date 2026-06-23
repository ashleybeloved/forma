package service

import (
	"forma/internal/model"
	"forma/internal/pkg"
	"forma/internal/repository"
)

type AuthService struct {
	Repo      *repository.UserRepository
	Validator *ValidatorService
}

func NewUserService(repo *repository.UserRepository, validator *ValidatorService) *AuthService {
	return &AuthService{
		Repo:      repo,
		Validator: validator,
	}
}

func (s *AuthService) Me(userID int) (*model.User, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Register(username, password string) (*model.User, error) {
	err := s.Validator.ValidateUsername(username)
	if err != nil {
		return nil, err
	}

	err = s.Validator.ValidatePassword(password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = s.Repo.AddUser(username, hashedPassword)
	if err != nil {
		return nil, err
	}

	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	user, err := s.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	compare := pkg.CompareHashAndPassword(user.Password, password)
	if !compare {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
