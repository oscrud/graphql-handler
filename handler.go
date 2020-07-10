package graphql

import (
	graphql "github.com/graphql-go/graphql"
	"github.com/oscrud/oscrud"
)

// Options :
type Options struct {
	RootObject          map[string]interface{}
	ReservedQueryString string
}

// Handler :
func Handler(schema graphql.Schema, opts ...Options) oscrud.Handler {
	options := Options{}
	if len(opts) > 0 {
		options = opts[0]
	}

	if options.ReservedQueryString == "" {
		options.ReservedQueryString = "query"
	}

	return func(ctx oscrud.Context) oscrud.Context {
		params := graphql.Params{
			Schema:     schema,
			RootObject: options.RootObject,
			Context:    ctx.Context(),
		}

		queries := ctx.Query()
		if graphQuery, ok := queries[options.ReservedQueryString]; ok {
			params.RequestString = graphQuery.(string)
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
