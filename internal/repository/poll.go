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

func NewPollRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (h *PollRepository) CreatePoll(poll *model.Poll, config []byte) error {
	result, err := h.DB.Exec(`INSERT INTO polls (title, description, config, creator_id, short_id) VALUES (?, ?, ?, ?, ?)`,
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

func (h *PollRepository) GetPollByID(pollID int) (*model.Poll, error) {
	poll := model.Poll{}

	err := h.DB.QueryRow(`SELECT id, title, description, config, creator_id, short_id, edited_at, created_at FROM polls WHERE id = ?`, pollID).
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

func (h *PollRepository) GetPollByShortID(pollShortID int) (*model.Poll, error) {
	poll := model.Poll{}

	err := h.DB.QueryRow(`SELECT id, title, description, config, creator_id, short_id edited_at, created_at FROM polls WHERE short_id = ?`, pollShortID).
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

func (h *PollRepository) GetPollsByCreatorID(creatorID int) (polls []model.Poll, err error) {
	rows, err := h.DB.Query(`SELECT (id, title, description, config, short_id, edited_at, created_at) FROM polls WHERE creator_id = ?`, creatorID)

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

		err := rows.Scan(&p.ID, &p.Title, &p.Description, &configBytes, &p.ShortID, &p.CreatedAt)
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

	return polls, nil
}
