package service

import (
	"forma/internal/config"
	"forma/internal/model"
	"regexp"
	"slices"
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

func (s *ValidatorService) ValidateAnswers(questions []model.Question, answers []model.Answer) error {
	if len(answers) != len(questions) {
		return ErrAnswersNotCompare
	}

	questionMap := make(map[int]model.Question, len(questions))
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for _, a := range answers {
		q, exists := questionMap[a.QuestionID]
		if !exists {
			return ErrQuestionNotFound
		}

		switch q.Type {
		case "text":
			if len(a.Options) != 1 {
				return ErrQuestionType
			}

			if a.Options[0] == "" {
				return ErrEmptyAnswer
			}
		case "single":
			if len(a.Options) != 1 {
				return ErrQuestionType
			}

			if !slices.Contains(q.Options, a.Options[0]) {
				return ErrOptionNotFound
			}
		case "multiply":
			if len(a.Options) > len(q.Options) || len(a.Options) < 1 {
				return ErrQuestionType
			}

			seenBefore := make(map[string]bool, len(a.Options))

			for _, o := range a.Options {
				if seenBefore[o] {
					return ErrDuplicateOptions
				}
				seenBefore[o] = true

				if !slices.Contains(q.Options, o) {
					return ErrOptionNotFound
				}
			}
		}
	}

	return nil
}
