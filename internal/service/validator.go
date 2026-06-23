package service

import (
	"forma/internal/config"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	hasLetter     = regexp.MustCompile(`[a-zA-Z]`)
)

type ValidatorService struct {
	Config *config.Config
}

func NewValidatorService(cfg *config.Config) *ValidatorService {
	return &ValidatorService{
		Config: cfg,
	}
}

func (s *ValidatorService) ValidateUsername(username string) error {
	usernameLength := utf8.RuneCountInString(username)

	if usernameLength > s.Config.UsernameMaxSymbols {
		return ErrUsernameTooBig
	}

	if usernameLength < s.Config.UsernameMinSymbols {
		return ErrUsernameTooSmall
	}

	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}

	if !hasLetter.MatchString(username) {
		return ErrMustContainLetter
	}

	if strings.Contains(username, "__") {
		return ErrDoubleUnderscores
	}

	return nil
}

func (s *ValidatorService) ValidatePassword(password string) error {
	passwordLength := utf8.RuneCountInString(password)

	if passwordLength > s.Config.PasswordMaxSymbols {
		return ErrPasswordTooBig
	}

	if passwordLength < s.Config.PasswordMinSymbols {
		return ErrPasswordTooSmall
	}

	return nil
}
