package repository

import (
	"database/sql"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

func Connect(dbPath string) *sql.DB {
	dbPathWAL := dbPath + "?_auth_journal_mode=WAL&_loc=auto"

	db, err := sql.Open("sqlite", dbPathWAL)
	if err != nil {
		slog.Error("failed to open database", "info:", err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		slog.Error("failed to connect database", "info:", err)
	}

	runMigrations(db)

	slog.Info("Database successfully loaded")

	return db
}

func runMigrations(db *sql.DB) {
	sqlFile, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		slog.Error("failed to read file 001_init.sql", "info:", err)
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		slog.Error("failed to execute migrations", "info:", err)
	}
}
