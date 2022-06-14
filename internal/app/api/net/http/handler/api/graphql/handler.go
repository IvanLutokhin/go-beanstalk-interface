package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/IvanLutokhin/go-beanstalk"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/executor"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/graphql/resolver"
	"net/http"
)

func Handler(pool beanstalk.Pool) http.Handler {
	server := handler.NewDefaultServer(executor.NewExecutableSchema(executor.Config{Resolvers: resolver.NewResolver(pool)}))

	server.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		if fields, ok := ctx.Value("middleware:logging:fields").(map[string]interface{}); ok {
			operationCtx := graphql.GetOperationContext(ctx)

			fields["graphql_operation"] = operationCtx.Operation.Operation
			fields["graphql_operation_name"] = operationCtx.Operation.Name
			fields["graphql_variables"] = operationCtx.Variables
		}

		return next(ctx)
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.ServeHTTP(w, r)
	})
}
