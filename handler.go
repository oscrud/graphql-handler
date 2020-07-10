package graphql

import (
	graphql "github.com/graphql-go/graphql"
	"github.com/oscrud/oscrud"
)

// Handler :
func Handler(operation string, reservedQueryString string, rootObject map[string]interface{}, schema graphql.Schema) oscrud.Handler {
	if reservedQueryString == "" {
		reservedQueryString = "query"
	}

	return func(ctx oscrud.Context) oscrud.Context {
		params := graphql.Params{
			Schema:        schema,
			OperationName: operation,
			RootObject:    rootObject,
			Context:       ctx.Context(),
		}

		queries := ctx.Query()
		if graphQuery, ok := queries[reservedQueryString]; ok {
			params.RequestString = graphQuery.(string)
			delete(queries, reservedQueryString)
			if len(queries) > 0 {
				params.VariableValues = queries
			}
		}

		result := graphql.Do(params)
		status := 200
		if len(result.Errors) > 0 {
			status = 400
		}
		return ctx.JSON(status, result).End()
	}
}
