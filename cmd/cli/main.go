package main

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/cli/handler"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/version"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

func main() {
	app := &cli.App{
		Name:    "beanstalk-cli",
		Usage:   "Provides Beanstalk queue commands",
		Version: version.Tag(),
		Commands: cli.Commands{
			{
				Name:      "put",
				Usage:     "puts job",
				ArgsUsage: "[data]",
				Action:    handler.CreateAction((*handler.Handler).Put),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.IntFlag{
						Name:  "priority",
						Value: 0,
						Usage: "jobs with smaller priority values will be scheduled before jobs with larger priorities",
					},
					&cli.DurationFlag{
						Name:  "delay",
						Value: 0,
						Usage: "the number of seconds to wait before putting the job in the ready queue",
					},
					&cli.DurationFlag{
						Name:  "ttr",
						Value: 0,
						Usage: "the number of seconds to allow a worker to run this job",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "delete",
				Usage:     "deletes job",
				ArgsUsage: "[job-id]",
				Action:    handler.CreateAction((*handler.Handler).Delete),
			},
			{
				Name:   "delete-jobs",
				Usage:  "deletes jobs",
				Action: handler.CreateAction((*handler.Handler).DeleteJobs),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.IntFlag{
						Name:  "count",
						Value: 1,
						Usage: "the number of jobs actually deleted",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "peek",
				Usage:     "peeks job",
				ArgsUsage: "[job-id]",
				Action:    handler.CreateAction((*handler.Handler).Peek),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:   "peek-ready",
				Usage:  "peeks ready job",
				Action: handler.CreateAction((*handler.Handler).PeekReady),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:   "peek-delayed",
				Usage:  "peeks delayed job",
				Action: handler.CreateAction((*handler.Handler).PeekDelayed),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:   "peek-buried",
				Usage:  "peeks buried job",
				Action: handler.CreateAction((*handler.Handler).PeekBuried),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "kick",
				Usage:     "kicks jobs",
				ArgsUsage: "[bound]",
				Action:    handler.CreateAction((*handler.Handler).Kick),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
				},
			},
			{
				Name:      "kick-job",
				Usage:     "kicks job",
				ArgsUsage: "[job-id]",
				Action:    handler.CreateAction((*handler.Handler).KickJob),
			},
			{
				Name:   "kick-jobs",
				Usage:  "kicks jobs at tube",
				Action: handler.CreateAction((*handler.Handler).KickJobs),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "tube",
						Value: "default",
						Usage: "the name of the tube now being used",
					},
					&cli.IntFlag{
						Name:  "count",
						Value: 1,
						Usage: "the number of jobs actually kicked",
					},
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "stats-job",
				Usage:     "gets job stats",
				ArgsUsage: "[job-id]",
				Action:    handler.CreateAction((*handler.Handler).StatsJob),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "stats-tube",
				Usage:     "gets tube stats",
				ArgsUsage: "[tube]",
				Action:    handler.CreateAction((*handler.Handler).StatsTube),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:   "stats",
				Usage:  "gets server stats",
				Action: handler.CreateAction((*handler.Handler).Stats),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:   "list-tubes",
				Usage:  "gets list of tubes",
				Action: handler.CreateAction((*handler.Handler).ListTubes),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "format",
						Value: "yaml",
						Usage: "the output format",
					},
				},
			},
			{
				Name:      "pause-tube",
				Usage:     "sets pause of tube",
				ArgsUsage: "[tube]",
				Action:    handler.CreateAction((*handler.Handler).PauseTube),
				Flags: []cli.Flag{
					&cli.DurationFlag{
						Name:  "delay",
						Value: 0,
						Usage: "the number of seconds to wait before reserving any more jobs from the queue",
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "quiet",
				Aliases: []string{"q"},
				Value:   false,
				Usage:   "do not output any message",
			},
			&cli.StringFlag{
				Name:    "host",
				EnvVars: []string{"BI_BEANSTALK_HOST"},
				Usage:   "the beanstalk host",
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   11300,
				EnvVars: []string{"BI_BEANSTALK_PORT"},
				Usage:   "the beanstalk port",
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("quiet") {
				ctx.App.Writer = io.Discard
				ctx.App.ErrWriter = io.Discard
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprint(os.Stderr, err)

		os.Exit(1)
	}
}
