CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS polls (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    config TEXT NOT NULL,
    creator_id INTEGER NOT NULL,
    short_id TEXT UNIQUE NOT NULL,
    edited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    poll_short_id TEXT NOT NULL,
    user_id INTEGER,
    ip TEXT NOT NULL,
    guest_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (poll_short_id) REFERENCES polls(short_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vote_answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vote_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    options TEXT NOT NULL,

    FOREIGN KEY (vote_id) REFERENCES votes(id) ON DELETE CASCADE
);

--

CREATE INDEX IF NOT EXISTS idx_votes_poll_short_id ON votes(poll_short_id);
CREATE INDEX IF NOT EXISTS idx_vote_answers_vote_id ON vote_answers(vote_id);
CREATE INDEX IF NOT EXISTS idx_vote_answers_analytics ON vote_answers(question_id, options);
