package graph

import (
	"cashflow/api/graph/generated"
	"cashflow/models"
	"context"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

type Directives struct {
	*models.App // App is a dependency
}

func GetDirectives(app *models.App) generated.DirectiveRoot {
	directives := &Directives{App: app}

	return generated.DirectiveRoot{
		Auth: directives.AuthDirective,
	}
}

func (d *Directives) AuthDirective(ctx context.Context, obj interface{}, next graphql.Resolver, requires *generated.Role) (interface{}, error) {
	if requires == nil {
		return next(ctx)
	}

	valid, err := (*d.AuthService).ValidateSession(ctx, requires.String())
	if err != nil || !valid {
		log.Println("access denied for role", requires.String())
		return nil, fmt.Errorf("access denied")
	}
	return next(ctx)
}
