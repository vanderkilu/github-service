package cmd

import (
	"log"
	"os"
	"time"

	"github.com/google/go-github/v63/github"
	"github.com/jackc/pgx/v4"
	"github.com/urfave/cli/v2"
	"github.com/vanderkilu/github-service/dao/postgresql"
	"github.com/vanderkilu/github-service/service"
)

func convertStringToTime(timeString string) (*time.Time) {
    layout := "2006-01-02T15:04:05Z"
    parsedTime, err := time.Parse(layout, timeString)
    if err != nil {
        return nil
    }
    return &parsedTime
}

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


			srv := service.New(postgresql.NewErrNoRowsQueries(postgresql.New(conn)), githubClient)

			time := convertStringToTime(os.Getenv("SINCE"))
			
			return srv.MonitorRepo(c.Context, os.Getenv("OWNER"), os.Getenv("REPO"), time)
		},
	}
	return &command
}