CREATE TABLE IF NOT EXISTS reviewers(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    pr_id INTEGER NOT NULL,
    first_reviewer_id INTEGER,
    second_reviewer_id INTEGER,

    UNIQUE(pr_id, first_reviewer_id),
    UNIQUE(pr_id, second_reviewer_id)
);