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

type NewPollRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config" binding:"required"`
}

type PollHeader struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config"`
	CreatorID   int        `json:"creator_id"`
	ShortID     string     `json:"short_id"`
	EditedAt    *time.Time `json:"edited_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
