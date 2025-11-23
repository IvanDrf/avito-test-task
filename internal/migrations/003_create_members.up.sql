CREATE TABLE IF NOT EXISTS members(
    team_id TEXT NOT NULL,
    user_id TEXT NOT NULL,

    FOREIGN KEY (team_id) REFERENCES teams(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE (team_id, user_id)
);
