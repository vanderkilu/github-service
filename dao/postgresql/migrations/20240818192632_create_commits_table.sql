-- +goose Up
CREATE TABLE commit 
(
    id SERIAL PRIMARY KEY,
    repo_full_name VARCHAR(255) NOT NULL,
    sha VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    url VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (repo_full_name) REFERENCES  repository(repo_full_name),
    UNIQUE(sha)
);

-- +goose Down
DROP TABLE commit;

