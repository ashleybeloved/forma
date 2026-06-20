package service

import (
	"forma/internal/model"
	"forma/internal/pkg"
	"forma/internal/repository"
)

type AuthService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (h *AuthService) Register(username, password string) (*model.User, error) {
	hashedPassword, err := pkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	err = h.Repo.AddUser(username, hashedPassword)
	if err != nil {
		return nil, err
	}

	user, err := h.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *AuthService) Login(username, password string) (*model.User, error) {
	user, err := h.Repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	compare := pkg.CompareHashAndPassword(user.Password, password)
	if !compare {
		return nil, ErrInvalidPassword
	}

	return user, nil
}
