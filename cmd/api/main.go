package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/log"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/net/http/middleware"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
	"github.com/IvanLutokhin/go-beanstalk-interface/pkg/version"
	"go.uber.org/fx"
	"os"
	"time"
)

func main() {
	configPath := flag.String("c", "", "Path to config file")

	flag.Parse()

	c, err := config.LoadOrDefault(*configPath)
	if err != nil {
		exitOnError(err)
	}

	app := fx.New(
		fx.Options(
			fx.NopLogger,
			fx.StartTimeout(5*time.Minute),
			fx.StopTimeout(5*time.Minute),
		),
		fx.Provide(
			func() *config.Config { return c },
			log.NewLogger,
			beanstalk.NewLoggerAdapter,
			beanstalk.NewPool,
			http.NewRouter,
			middleware.NewCors,
			middleware.NewLogging,
			middleware.NewRecovery,
			security.NewUserProvider,
		),
		fx.Invoke(
			registerHooks,
			security.RegisterHooks,
			beanstalk.RegisterHooks,
			http.RegisterHooks,
		),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()

	if err = app.Start(startCtx); err != nil {
		exitOnError(err)
	}

	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()

	if err = app.Stop(stopCtx); err != nil {
		exitOnError(err)
	}
}

func registerHooks(lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println(version.String())

			return nil
		},
	})
}

func exitOnError(err error) {
	fmt.Fprint(os.Stderr, err)

	os.Exit(1)
}
