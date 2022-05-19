package api

import (
	"context"
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/log"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/version"
	"go.uber.org/fx"
	"time"
)

func New() *fx.App {
	return fx.New(
		fx.Options(
			fx.NopLogger,
			fx.StartTimeout(5*time.Minute),
			fx.StopTimeout(5*time.Minute),
		),
		fx.Provide(
			config.New,
			log.NewLogger,
			beanstalk.NewPool,
			http.NewRouter,
			middleware.NewCors,
			middleware.NewLogging,
			middleware.NewRecovery,
			security.NewUserProvider,
		),
		fx.Invoke(
			RegisterHooks,
			security.RegisterHooks,
			beanstalk.RegisterHooks,
			http.RegisterHooks,
		),
	)
}

func RegisterHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println(version.String())

			return nil
		},
	})
}
