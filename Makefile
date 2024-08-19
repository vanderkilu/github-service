build:
	go build -o bin/github-service ./main.go
run:
	DATABASE_URL=postgres://root:root@localhost:5455/perfectscale?search_path=github_service OWNER=chromium REPO=chromium  ./bin/github-service cron-monitor-repo