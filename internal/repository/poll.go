package repository

import (
	"database/sql"
	"encoding/json"
	"forma/internal/model"
	"log/slog"
)

type PollRepository struct {
	DB *sql.DB
}

func NewPollRepository(db *sql.DB) *PollRepository {
	return &PollRepository{
		DB: db,
	}
}

func (r *PollRepository) CreatePoll(poll *model.Poll, config []byte) error {
	result, err := r.DB.Exec(`INSERT INTO polls (title, description, config, creator_id, short_id) VALUES (?, ?, ?, ?, ?)`,
		poll.Title, poll.Description, config, poll.CreatorID, poll.ShortID)

	if err != nil {
		slog.Error("failed to execute query", "error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("failed to insert last id to poll structure", "error", err)
		return err
	}

	poll.ID = int(id)

	return nil
}

func (r *PollRepository) UpdatePoll(poll *model.Poll, configBytes []byte, userID int) error {
	result, err := r.DB.Exec(`UPDATE polls SET title = ?, description = ?, config = ?, edited_at = DATETIME('now') WHERE id = ? AND creator_id = ?`,
		poll.Title, poll.Description, configBytes, poll.ID, userID)

	if err != nil {
		slog.Error("failed to execute query", "error", err)
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrPollNotFound
	}

	return nil
}

func (r *PollRepository) DeletePoll(pollID, creatorID int) error {
	result, err := r.DB.Exec(`DELETE FROM polls WHERE id = ? AND creator_id = ?`, pollID, creatorID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrPollNotFound
	}

	return nil
}

func (r *PollRepository) GetPollByID(pollID int) (*model.Poll, error) {
	poll := model.Poll{}

	err := r.DB.QueryRow(`SELECT id, title, description, config, creator_id, short_id, edited_at, created_at FROM polls WHERE id = ?`, pollID).
		Scan(&poll.ID, &poll.Title, &poll.Description, &poll.Config, &poll.CreatorID, &poll.ShortID, &poll.EditedAt, &poll.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPollNotFound
		}

		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	return &poll, nil
}

func (r *PollRepository) GetPollByShortID(pollShortID string) (*model.Poll, error) {
	var poll = model.Poll{}
	var configBytes []byte

	err := r.DB.QueryRow(`SELECT id, title, description, config, creator_id, short_id, edited_at, created_at FROM polls WHERE short_id = ?`, pollShortID).
		Scan(&poll.ID, &poll.Title, &poll.Description, &configBytes, &poll.CreatorID, &poll.ShortID, &poll.EditedAt, &poll.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPollNotFound
		}

		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	err = json.Unmarshal(configBytes, &poll.Config)
	if err != nil {
		return nil, err
	}

	return &poll, nil
}

func (r *PollRepository) GetPollsByCreatorID(creatorID, limit, offset int) (polls []model.Poll, err error) {
	rows, err := r.DB.Query(`SELECT id, title, description, config, short_id, edited_at, created_at FROM polls WHERE creator_id = ? LIMIT ? OFFSET ?`, creatorID, limit, offset)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPollsNotFound
		}

		slog.Error("failed to execute query", "error", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var p model.Poll
		var configBytes []byte

		err := rows.Scan(&p.ID, &p.Title, &p.Description, &configBytes, &p.ShortID, &p.EditedAt, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(configBytes, &p.Config)
		if err != nil {
			slog.Warn("failed to unmarshal json from polls", "error", err)
			return nil, err
		}

		p.CreatorID = creatorID

		polls = append(polls, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if polls == nil {
		return nil, ErrPollsNotFound
	}

	return polls, nil
}

func (r *PollRepository) Vote(vote *model.Vote, answers []model.Answer) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	result, err := tx.Exec(`INSERT INTO votes (poll_short_id, user_id, ip, country_code, guest_token) VALUES (?, ?, ?, ?, ?)`,
		vote.PollShortID, vote.UserID, vote.IP, vote.CountryCode, vote.GuestToken)

	if err != nil {
		slog.Error("failed to execute query", "error", err)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		slog.Error("failed to insert last id to vote structure", "error", err)
		return err
	}

	for _, answer := range answers {
		for _, option := range answer.Options {
			_, err := tx.Exec(`INSERT INTO vote_answers (vote_id, question_id, options) VALUES (?, ?, ?)`,
				int(id), answer.QuestionID, option)

			if err != nil {
				slog.Error("failed to insert vote answer",
					"error", err,
					"vote_id", int(id),
					"question_id", answer.QuestionID,
				)
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit vote transaction", "error", err)
		return err
	}

	return nil
}

func (r *PollRepository) HasVoted(secured bool, pollShortID string, ip string, guestToken string, userID int) (bool, error) {
	querySecured := `
		SELECT EXISTS(
			SELECT 1 FROM votes
			WHERE poll_short_id = ?
			  AND (
				guest_token = ?
				OR ip = ?
				OR (user_id = ? AND user_id != -1)
			  )
		);`

	query := `SELECT EXISTS(
		SELECT 1 FROM votes
		WHERE poll_short_id = ?
		  AND (
			guest_token = ?
			OR (user_id = ? AND user_id != -1)
		  )
	);`

	var exists bool
	var err error

	if secured {
		err = r.DB.QueryRow(querySecured, pollShortID, guestToken, ip, userID).Scan(&exists)
	} else {
		err = r.DB.QueryRow(query, pollShortID, guestToken, userID).Scan(&exists)
	}

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PollRepository) PollStatistics(shortID string) error {
	return nil
}
