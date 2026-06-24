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
	Secured     bool       `json:"secured"`
	AuthOnly    bool       `json:"auth_only"`
	EditedAt    time.Time  `json:"edited_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type NewPollRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config" binding:"required"`
	Secured     *bool      `json:"secured" binding:"required"`
	AuthOnly    *bool      `json:"auth_only" binding:"required"`
}

type UpdatePollRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Config      PollConfig `json:"config" binding:"required"`
	Secured     *bool      `json:"secured" binding:"required"`
	AuthOnly    *bool      `json:"auth_only" binding:"required"`
}

type Answer struct {
	QuestionID int      `json:"question_id" binding:"required"`
	Options    []string `json:"options" binding:"required"`
}

type Vote struct {
	ID          int      `json:"id"`
	PollShortID string   `json:"poll_short_id"`
	UserID      int      `json:"user_id"`
	IP          string   `json:"ip"`
	CountryCode string   `json:"country_code"`
	GuestToken  string   `json:"guest_token"`
	Answers     []Answer `json:"answers"`
}

type NewVoteRequest struct {
	Answers []Answer `json:"answers" binding:"required"`
}

type Stats struct {
	TotalVotes      int              `json:"total_votes"`
	QuestionResults []QuestionResult `json:"results"`
	TopCountries    []CountryResult  `json:"top_countries"`
}

type QuestionResult struct {
	QuestionID int      `json:"id"`
	Options    []Result `json:"options"`
}

type Result struct {
	Option     string  `json:"option"`
	Votes      int     `json:"votes"`
	Percentage float64 `json:"percentage"`
}

type CountryResult struct {
	CountryCode string `json:"country_code"`
	Votes       int    `json:"votes"`
}
