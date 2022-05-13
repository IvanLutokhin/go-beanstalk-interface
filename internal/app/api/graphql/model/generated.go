// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"github.com/IvanLutokhin/go-beanstalk"
)

type BuryJobInput struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
}

type BuryJobPayload struct {
	ID int `json:"id"`
}

type CreateJobInput struct {
	Tube     string `json:"tube"`
	Priority int    `json:"priority"`
	Delay    int    `json:"delay"`
	Ttr      int    `json:"ttr"`
	Data     string `json:"data"`
}

type CreateJobPayload struct {
	Tube string `json:"tube"`
	ID   int    `json:"id"`
}

type DeleteJobInput struct {
	ID int `json:"id"`
}

type DeleteJobPayload struct {
	ID int `json:"id"`
}

type Job struct {
	ID    int                 `json:"id"`
	Data  string              `json:"data"`
	Stats *beanstalk.StatsJob `json:"stats"`
}

type KickJobInput struct {
	ID int `json:"id"`
}

type KickJobPayload struct {
	ID int `json:"id"`
}

type ReleaseJobInput struct {
	ID       int `json:"id"`
	Priority int `json:"priority"`
	Delay    int `json:"delay"`
}

type ReleaseJobPayload struct {
	ID int `json:"id"`
}

type Server struct {
	Stats *beanstalk.Stats `json:"stats"`
}

type Tube struct {
	Name       string               `json:"name"`
	Stats      *beanstalk.StatsTube `json:"stats"`
	ReadyJob   *Job                 `json:"readyJob"`
	DelayedJob *Job                 `json:"delayedJob"`
	BuriedJob  *Job                 `json:"buriedJob"`
}

type TubeConnection struct {
	Edges []TubeEdge `json:"edges"`
}

type TubeEdge struct {
	Node *Tube `json:"node"`
}
