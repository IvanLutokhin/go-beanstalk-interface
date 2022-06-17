package v1

import (
	"encoding/json"
	"errors"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

func GetEmbedFiles(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Clean(r.URL.Path)

		if strings.HasSuffix(p, "/") {
			writer.JSON(w, http.StatusNotFound, response.NotFound())

			return
		}

		_, err := fs.Open(p)
		if os.IsNotExist(err) {
			writer.JSON(w, http.StatusNotFound, response.NotFound())

			return
		}

		http.FileServer(fs).ServeHTTP(w, r)
	})
}

// swagger:route GET /server/stats server server-stats
//
// Beanstalk server statistics
//
// Responses:
// 200: ServerStatsSuccessResponse
//
// Security:
// - basicAuth: []
func GetServerStats() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		stats, err := c.Stats()
		if err != nil {
			panic(err)
		}

		writer.JSON(w, http.StatusOK, response.Success(ServerStatsResponse{stats}))
	})
}

// swagger:route GET /tubes tubes list-of-tubes
//
// List of exists tubes
//
// Responses:
// 200: TubesSuccessResponse
//
// Security:
// - basicAuth: []
func GetTubes() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		tubes, err := c.ListTubes()
		if err != nil {
			panic(err)
		}

		writer.JSON(w, http.StatusOK, response.Success(TubesResponse{Tubes: tubes}))
	})
}

// swagger:route GET /tubes/{name}/stats tubes tube-stats
//
// Gets statistics of specified tube
//
// Responses:
// 200: TubeStatsSuccessResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func GetTubeStats() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(client beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		stats, err := client.StatsTube(vars["name"])
		if err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(TubeStatsResponse{stats}))
	})
}

// swagger:route POST /jobs jobs create-job
//
// Creates a new job
//
// Responses:
// 200: CreateJobSuccessResponse
// 400: BadRequestResponse
//
// Security:
// - basicAuth: []
func CreateJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		var request CreateJobRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			writer.JSON(w, http.StatusBadRequest, response.BadRequest(ErrorData{[]string{"failed to unmarshal request"}}))

			return
		}

		tube, err := c.Use(request.Tube)
		if err != nil {
			panic(err)
		}

		id, err := c.Put(request.Priority, time.Duration(request.Delay), time.Duration(request.TTR), []byte(request.Data))
		if err != nil {
			panic(err)
		}

		writer.JSON(w, http.StatusCreated, response.Success(CreateJobResponse{tube, id}))
	})
}

// swagger:route GET /jobs/{id} jobs get-job
//
// Gets data of the specified job
//
// Responses:
// 200: GetJobSuccessResponse
// 400: BadRequestResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func GetJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		job, err := c.Peek(id)
		if err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(JobResponse{string(job.Data)}))
	})
}

// swagger:route POST /jobs/{id}/bury jobs bury-job
//
// Bury the specified job
//
// Responses:
// 200: SuccessResponse
// 400: BadRequestResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func BuryJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var request BuryJobRequest
		if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
			writer.JSON(w, http.StatusBadRequest, response.BadRequest(ErrorData{[]string{"failed to unmarshal request"}}))

			return
		}

		if err = c.Bury(id, request.Priority); err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(nil))
	})
}

// swagger:route POST /jobs/{id}/delete jobs delete-job
//
// Delete the specified job
//
// Responses:
// 200: SuccessResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func DeleteJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		if err = c.Delete(id); err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(nil))
	})
}

// swagger:route POST /jobs/{id}/kick jobs kick-job
//
// Kick the specified job
//
// Responses:
// 200: SuccessResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func KickJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		if err = c.KickJob(id); err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(nil))
	})
}

// swagger:route POST /jobs/{id}/release jobs release-job
//
// Release the specified job
//
// Responses:
// 200: SuccessResponse
// 400: BadRequestResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func ReleaseJob() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		var request ReleaseJobRequest
		if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
			writer.JSON(w, http.StatusBadRequest, response.BadRequest(ErrorData{[]string{"failed to unmarshal request"}}))

			return
		}

		if err = c.Release(id, request.Priority, time.Duration(request.Delay)); err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(nil))
	})
}

// swagger:route GET /jobs/{id}/stats jobs job-stats
//
// Statistics of specified job
//
// Responses:
// 200: JobStatsSuccessResponse
// 404: NotFoundResponse
//
// Security:
// - basicAuth: []
func GetJobStats() beanstalk.Handler {
	return beanstalk.HandlerFunc(func(c beanstalk.Client, w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			panic(err)
		}

		stats, err := c.StatsJob(id)
		if err != nil {
			if errors.Is(err, beanstalk.ErrNotFound) {
				writer.JSON(w, http.StatusNotFound, response.NotFound())

				return
			} else {
				panic(err)
			}
		}

		writer.JSON(w, http.StatusOK, response.Success(JobStatsResponse{stats}))
	})
}
