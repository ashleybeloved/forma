package model

import "time"

type Question struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Options     []string `json:"options"`
}

type PollConfig struct {
	Questions []Question `json:"questions"`
}

type Poll struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config"`
	CreatorID   int        `json:"creator_id"`
	ShortID     string     `json:"short_id"`
	EditedAt    time.Time  `json:"edited_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Answer struct {
	QuestionID int      `json:"question_id" binding:"required"`
	Options    []string `json:"options" binding:"required"`
}

type Answers struct {
	Answers []Answer `json:"answers" binding:"required"`
}

type Vote struct {
	ID          int     `json:"id"`
	PollShortID string  `json:"poll_short_id"`
	UserID      int     `json:"user_id"`
	IP          string  `json:"ip"`
	GuestToken  string  `json:"guest_token"`
	Answers     Answers `json:"answers"`
}

type NewVoteRequest struct {
	Answers Answers `json:"answers" binding:"required"`
}

type NewPollRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config" binding:"required"`
}

type UpdatePollRequest struct {
	ID          int        `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config" binding:"required"`
}

type DeletePollRequest struct {
	ID int `json:"id" binding:"required"`
}
