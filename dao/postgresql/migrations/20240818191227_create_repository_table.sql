-- +goose Up
CREATE TABLE repository 
(
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    url VARCHAR(255) NOT NULL,
    language VARCHAR(100) NOT NULL,
    repo_name VARCHAR(255) NOT NULL,
    repo_full_name VARCHAR(255) NOT NULL,
    forks_count INT DEFAULT 0,
    stars_count INT DEFAULT 0,
    open_issues_count INT DEFAULT 0,
    watchers_count INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (repo_full_name)
);

-- +goose Down
DROP TABLE repository;

