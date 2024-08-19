-- name: CreateRepository :exec
INSERT INTO repository (id, description, url, language,repo_name, repo_full_name, forks_count, stars_count, open_issues_count, watchers_count,created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (repo_full_name) DO NOTHING;

-- name: CreateCommit :exec
INSERT INTO commit (id, repo_full_name, sha, message, url, author, date) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (sha) DO NOTHING;

-- name: GetLastCommitSha :one
SELECT sha FROM commit ORDER BY date DESC LIMIT 1;