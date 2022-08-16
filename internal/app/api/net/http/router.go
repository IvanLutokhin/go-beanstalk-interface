package http

import (
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/api"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/graphql"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/api/system/v1"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/handler/auth"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/response"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/writer"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/embed"
	"github.com/gorilla/mux"
	"io/fs"
	"net/http"
	"path"
)

func NewRouter(
	config *config.Config,
	cors *middleware.Cors,
	logging *middleware.Logging,
	recovery *middleware.Recovery,
	provider *security.UserProvider,
	manager *security.TokenManager,
	pool beanstalk.Pool,
) *mux.Router {
	router := mux.NewRouter()

	router.StrictSlash(true)

	router.Use(
		logging.Middleware,
		recovery.Middleware,
		cors.Middleware,
	)

	registerAuthRoutes(router, config, provider, manager)

	registerAPIRoutes(router, provider, manager, pool)

	return router
}

func registerAuthRoutes(router *mux.Router, config *config.Config, provider *security.UserProvider, manager *security.TokenManager) {
	sr := router.PathPrefix("/auth").Subrouter()

	sr.Methods(http.MethodOptions, http.MethodPost).Path("/token").Handler(auth.Token(config, provider, manager))
	sr.Methods(http.MethodOptions, http.MethodPost).Path("/login").Handler(auth.Login(config, provider, manager))
	sr.Methods(http.MethodOptions, http.MethodPost).Path("/logout").Handler(auth.Logout())
}

func registerAPIRoutes(router *mux.Router, provider *security.UserProvider, manager *security.TokenManager, pool beanstalk.Pool) {
	sr := router.PathPrefix("/api").Subrouter()

	sr.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusNotFound, response.NotFound())
	})

	sr.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer.JSON(w, http.StatusMethodNotAllowed, response.MethodNotAllowed())
	})

	registerSystemV1Routes(sr, provider, manager, pool)

	registerGraphQLRoutes(sr, provider, manager, pool)
}

func registerSystemV1Routes(router *mux.Router, provider *security.UserProvider, manager *security.TokenManager, pool beanstalk.Pool) {
	h := func(scopes []security.Scope, handler beanstalk.Handler) http.Handler {
		return middleware.Auth(provider, manager, scopes).Middleware(beanstalk.NewHTTPHandlerAdapter(pool, handler))
	}

	sr := router.PathPrefix("/system/v1").Subrouter()

	sr.Methods(http.MethodGet).Path("/server/stats").Handler(h([]security.Scope{security.ScopeReadServer}, v1.GetServerStats()))
	sr.Methods(http.MethodGet).Path("/tubes").Handler(h([]security.Scope{security.ScopeReadTubes}, v1.GetTubes()))
	sr.Methods(http.MethodGet).Path("/tubes/{name}/stats").Handler(h([]security.Scope{security.ScopeReadTubes}, v1.GetTubeStats()))
	sr.Methods(http.MethodPost).Path("/jobs").Handler(h([]security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}, v1.CreateJob()))
	sr.Methods(http.MethodGet).Path("/jobs/{id:[0-9]+}").Handler(h([]security.Scope{security.ScopeReadJobs}, v1.GetJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/bury").Handler(h([]security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}, v1.BuryJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/delete").Handler(h([]security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}, v1.DeleteJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/kick").Handler(h([]security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}, v1.KickJob()))
	sr.Methods(http.MethodPost).Path("/jobs/{id:[0-9]+}/release").Handler(h([]security.Scope{security.ScopeReadJobs, security.ScopeWriteJobs}, v1.ReleaseJob()))
	sr.Methods(http.MethodGet).Path("/jobs/{id:[0-9]+}/stats").Handler(h([]security.Scope{security.ScopeReadJobs}, v1.GetJobStats()))

	sr.PathPrefix("/").Handler(http.StripPrefix("/api/system/v1", v1.GetEmbedFiles(http.FS(embed.FSFunc(func(name string) (fs.File, error) {
		return api.SystemV1EmbedFS.Open(path.Join("system/v1", name))
	})))))
}

func registerGraphQLRoutes(router *mux.Router, provider *security.UserProvider, manager *security.TokenManager, pool beanstalk.Pool) {
	sr := router.PathPrefix("/graphql").Subrouter()

	sr.
		Methods(http.MethodGet, http.MethodPost).
		Path("").
		Handler(middleware.Auth(provider, manager, []security.Scope{}).Middleware(graphql.Handler(pool)))
}
