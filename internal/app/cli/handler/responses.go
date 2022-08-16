package handler

type JobResponse struct {
	ID   int    `json:"id" yaml:"id"`
	Data string `json:"data" yaml:"data"`
}

type PutCommandResponse struct {
	Tube string `json:"tube" yaml:"tube"`
	ID   int    `json:"id" yaml:"id"`
}

type DeleteJobsCommandResponse struct {
	Tube  string `json:"tube" yaml:"tube"`
	Count int    `json:"count" yaml:"count"`
}

type KickCommandResponse struct {
	Tube  string `json:"tube" yaml:"tube"`
	Count int    `json:"count" yaml:"count"`
}

type KickJobsCommandResponse struct {
	Tube  string `json:"tube" yaml:"tube"`
	Count int    `json:"count" yaml:"count"`
}
