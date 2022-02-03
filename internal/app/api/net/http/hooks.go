package http

import (
	"context"
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/gorilla/mux"
	"go.uber.org/fx"
	"net/http"
	"time"
)

func RegisterHooks(lifecycle fx.Lifecycle, config *config.Config, router *mux.Router) {
	server := &http.Server{
		Addr:         config.Http.ListenAddresses,
		ReadTimeout:  config.Http.ReadTimeout * time.Second,
		WriteTimeout: config.Http.WriteTimeout * time.Second,
		IdleTimeout:  config.Http.IdleTimeout * time.Second,
		Handler:      router,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
					panic(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
