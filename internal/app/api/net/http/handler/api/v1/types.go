package v1

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"time"
)

const (
	JobDataFormatRAW    = "raw"
	JobDataFormatBase64 = "base64"
)

type ErrorData struct {
	Errors []string `json:"errors"`
}

type JobData struct {
	Format  string `json:"format"`
	Payload string `json:"payload"`
}

// Request types
type (
	CreateJobRequest struct {
		Tube     string        `json:"tube"`
		Priority uint32        `json:"priority"`
		Delay    time.Duration `json:"delay"`
		TTR      time.Duration `json:"ttr"`
		Data     JobData       `json:"data"`
	}

	BuryJobRequest struct {
		Priority uint32 `json:"priority"`
	}

	ReleaseJobRequest struct {
		Priority uint32        `json:"priority"`
		Delay    time.Duration `json:"delay"`
	}
)

// Response types
type (
	ServerStatsResponse struct {
		Stats beanstalk.Stats `json:"stats"`
	}

	TubesResponse struct {
		Tubes []string `json:"tubes"`
	}

	TubeStatsResponse struct {
		Stats beanstalk.StatsTube `json:"stats"`
	}

	CreateJobResponse struct {
		Tube string `json:"tube"`
		ID   int    `json:"id"`
	}

	JobResponse struct {
		Data JobData `json:"data"`
	}

	JobStatsResponse struct {
		Stats beanstalk.StatsJob `json:"stats"`
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
type SwaggerTubeStatsResponse struct {
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
	// Job data format
	//
	// in: query
	// example: raw
	Format string `json:"format"`
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
