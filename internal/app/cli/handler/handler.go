package handler

import (
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/cli/printer"
	"github.com/urfave/cli/v2"
	"strconv"
)

type ActionFunc func(*Handler, *cli.Context) error

func CreateAction(f ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		client, err := beanstalk.Dial(fmt.Sprintf("%s:%d", ctx.String("host"), ctx.Int("port")))
		if err != nil {
			return err
		}

		defer client.Close()

		return f(&Handler{Client: client}, ctx)
	}
}

type Handler struct {
	Client beanstalk.Client
}

func (h *Handler) Put(ctx *cli.Context) error {
	data := ctx.Args().First()
	if len(data) == 0 {
		return RequiredArgumentError("data")
	}

	tube, err := h.Client.Use(ctx.String("tube"))
	if err != nil {
		return err
	}

	id, err := h.Client.Put(uint32(ctx.Int("priority")), ctx.Duration("delay"), ctx.Duration("ttr"), []byte(data))
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, PutCommandResponse{Tube: tube, ID: id})
}

func (h *Handler) Delete(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	return h.Client.Delete(jobID)
}

func (h *Handler) Release(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	return h.Client.Release(jobID, uint32(ctx.Int("priority")), ctx.Duration("delay"))
}

func (h *Handler) Bury(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	return h.Client.Bury(jobID, uint32(ctx.Int("priority")))
}

func (h *Handler) Peek(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	job, err := h.Client.Peek(jobID)
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, JobResponse{ID: job.ID, Data: string(job.Data)})
}

func (h *Handler) PeekReady(ctx *cli.Context) error {
	_, err := h.Client.Use(ctx.String("tube"))
	if err != nil {
		return err
	}

	job, err := h.Client.PeekReady()
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, JobResponse{ID: job.ID, Data: string(job.Data)})
}

func (h *Handler) PeekDelayed(ctx *cli.Context) error {
	_, err := h.Client.Use(ctx.String("tube"))
	if err != nil {
		return err
	}

	job, err := h.Client.PeekDelayed()
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, JobResponse{ID: job.ID, Data: string(job.Data)})
}

func (h *Handler) PeekBuried(ctx *cli.Context) error {
	_, err := h.Client.Use(ctx.String("tube"))
	if err != nil {
		return err
	}

	job, err := h.Client.PeekBuried()
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, JobResponse{ID: job.ID, Data: string(job.Data)})
}

func (h *Handler) Kick(ctx *cli.Context) error {
	bound, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("bound")
	}

	tube, err := h.Client.Use(ctx.String("tube"))
	if err != nil {
		return err
	}

	count, err := h.Client.Kick(bound)
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, KickCommandResponse{Tube: tube, Count: count})
}

func (h *Handler) KickJob(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	return h.Client.KickJob(jobID)
}

func (h *Handler) StatsJob(ctx *cli.Context) error {
	jobID, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		return InvalidArgumentError("job-id")
	}

	stats, err := h.Client.StatsJob(jobID)
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, stats)
}

func (h *Handler) StatsTube(ctx *cli.Context) error {
	tube := ctx.Args().First()
	if len(tube) == 0 {
		return RequiredArgumentError("tube")
	}

	stats, err := h.Client.StatsTube(tube)
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, stats)
}

func (h *Handler) Stats(ctx *cli.Context) error {
	stats, err := h.Client.Stats()
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, stats)
}

func (h *Handler) ListTubes(ctx *cli.Context) error {
	tubes, err := h.Client.ListTubes()
	if err != nil {
		return err
	}

	return printer.Print(ctx.String("format"), ctx.App.Writer, tubes)
}

func (h *Handler) PauseTube(ctx *cli.Context) error {
	tube := ctx.Args().First()
	if len(tube) == 0 {
		return RequiredArgumentError("tube")
	}

	return h.Client.PauseTube(tube, ctx.Duration("delay"))
}
