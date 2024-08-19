package cmd

import (
	"log"
	"os"

	"github.com/google/go-github/v63/github"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
	"github.com/vanderkilu/github-service/dao/postgresql"
	"github.com/vanderkilu/github-service/service"
)



func CronMonitorRepo() *cli.Command {
	command := cli.Command{
		Name:  "cron-monitor-repo",
		Usage: "cli cron to monitor and run repo",
		Action: func(c *cli.Context) error {

			conn, err := pgx.Connect(c.Context, os.Getenv("DATABASE_URL"))
			if err != nil {
				log.Fatalf("Failed to connect to the database: %v", err)
			}
			defer conn.Close(c.Context)

			githubClient := github.NewClient(nil)


			srv := service.New(postgresql.New(conn), githubClient)

			return srv.MonitorRepo(c.Context, os.Getenv("OWNER"), os.Getenv("REPO"), nil)
		},
	}
	return &command
}