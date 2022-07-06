package v1

import (
	"encoding/json"
	"errors"
	"github.com/IvanLutokhin/go-beanstalk"
	"time"
)

type ErrorData struct {
	Errors []string `json:"errors"`
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		*d = Duration(value)

		return nil

	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}

		*d = Duration(tmp)

		return nil

	default:
		return errors.New("invalid duration")
	}
}

// Request types
type (
	CreateJobRequest struct {
		Tube     string   `json:"tube"`
		Priority uint32   `json:"priority"`
		Delay    Duration `json:"delay"`
		TTR      Duration `json:"ttr"`
		Data     string   `json:"data"`
	}

	BuryJobRequest struct {
		Priority uint32 `json:"priority"`
	}

	ReleaseJobRequest struct {
		Priority uint32   `json:"priority"`
		Delay    Duration `json:"delay"`
	}
)

// Response types
type (
	ServerStatsResponse struct {
		Stats *beanstalk.Stats `json:"stats"`
	}

	TubesResponse struct {
		Tubes []string `json:"tubes"`
	}

	TubeStatsResponse struct {
		Stats *beanstalk.StatsTube `json:"stats"`
	}

	CreateJobResponse struct {
		Tube string `json:"tube"`
		ID   int    `json:"id"`
	}

	JobResponse struct {
		Data string `json:"data"`
	}

	JobStatsResponse struct {
		Stats *beanstalk.StatsJob `json:"stats"`
	}
)

// Swagger types

// swagger:response ServerStatsSuccessResponse
type ServerStatsSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string              `json:"status"`
		Data   ServerStatsResponse `json:"data"`
	}
}

// swagger:response TubesSuccessResponse
type TubesSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string        `json:"status"`
		Data   TubesResponse `json:"data"`
	}
}

// swagger:parameters tube-stats
type TubeStatsSwaggerParameters struct {
	// Tube name
	//
	// in: path
	// example: default
	Name string `json:"name"`
}

// swagger:response TubeStatsSuccessResponse
type TubeStatsSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string            `json:"status"`
		Data   TubeStatsResponse `json:"data"`
	}
}

// swagger:parameters create-job
type CreateJobSwaggerParameters struct {
	// in: body
	Body CreateJobRequest
}

// swagger:response CreateJobSuccessResponse
type CreateJobSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string            `json:"status"`
		Data   CreateJobResponse `json:"data"`
	}
}

// swagger:parameters get-job
type GetJobSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
}

// swagger:response GetJobSuccessResponse
type GetJobSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string      `json:"status"`
		Data   JobResponse `json:"data"`
	}
}

// swagger:parameters bury-job
type BuryJobSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
	// in: body
	Body BuryJobRequest
}

// swagger:parameters delete-job
type DeleteJobSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
}

// swagger:parameters kick-job
type KickJobSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
}

// swagger:parameters release-job
type ReleaseJobSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
	// in: body
	Body ReleaseJobRequest
}

// swagger:parameters job-stats
type JobStatsSwaggerParameters struct {
	// Job identifier
	//
	// in: path
	// example: 1
	ID int `json:"id"`
}

// swagger:response JobStatsSuccessResponse
type JobStatsSuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string           `json:"status"`
		Data   JobStatsResponse `json:"data"`
	}
}

// swagger:response SuccessResponse
type SuccessSwaggerResponse struct {
	// in: body
	Response struct {
		// default: success
		Status string `json:"status"`
	}
}

// swagger:response BadRequestResponse
type BadRequestSwaggerResponse struct {
	// in: body
	Response struct {
		// default: failure
		Status string `json:"status"`
		// default: Bad Request
		Message string    `json:"message"`
		Data    ErrorData `json:"data"`
	}
}

// swagger:response NotFoundResponse
type NotFoundSwaggerResponse struct {
	// in: body
	Response struct {
		// default: failure
		Status string `json:"status"`
		// default: Not Found
		Message string `json:"message"`
	}
}
