package graphql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/resolver"
	"net/http"
)

func Handler(pool beanstalk.Pool) http.Handler {
	return handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver.NewResolver(pool)}))
}
