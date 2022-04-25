package http

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/api"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/v1"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/embed"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
	"path"
)

func NewRouter(
	cors *middleware.Cors,
	logging *middleware.Logging,
	recovery *middleware.Recovery,
	pool beanstalk.Pool,
) *mux.Router {
	router := mux.NewRouter()

	router.StrictSlash(true)

	router.Use(
		recovery.Middleware,
		logging.Middleware,
		cors.Middleware,
	)

	registerAPIRoutes(router, pool)

	return router
}

func registerAPIRoutes(router *mux.Router, pool beanstalk.Pool) {
	sr := router.PathPrefix("/api/v1").Subrouter()

	sr.Methods(http.MethodGet).Path("/server/stats").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.GetServerStats()))
	sr.Methods(http.MethodGet).Path("/tubes").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.GetTubes()))
	sr.Methods(http.MethodGet).Path("/tubes/{name}/stats").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.GetTubeStats()))
	sr.Methods(http.MethodPost).Path("/jobs").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.CreateJob()))
	sr.Methods(http.MethodGet).Path("/jobs/{id:[0-9]+}").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/bury").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.BuryJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/delete").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.DeleteJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/kick").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.KickJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/release").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.ReleaseJob()))
	sr.Methods(http.MethodGet).Path("/jobs/{id:[0-9]+}/stats").Handler(beanstalk.NewHTTPHandlerAdapter(pool, v1.GetJobStats()))

	sr.PathPrefix("/").Handler(http.StripPrefix("/api/v1", v1.GetEmbedFiles(http.FS(embed.FSFunc(func(name string) (fs.File, error) {
		return api.V1EmbedFS.Open(path.Join("v1", name))
	})))))

	sr.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusNotFound, response.NotFound())
	})

	sr.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusMethodNotAllowed, response.MethodNotAllowed())
	})
}
