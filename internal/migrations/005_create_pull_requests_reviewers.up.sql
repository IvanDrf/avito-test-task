CREATE TABLE IF NOT EXISTS reviewers(
    id TEXT PRIMARY KEY,
    
    pr_id TEXT NOT NULL,
    first_reviewer_id TEXT,
    second_reviewer_id TEXT,

    UNIQUE(pr_id, first_reviewer_id),
    UNIQUE(pr_id, second_reviewer_id)
);