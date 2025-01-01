package graph

import (
	"cashflow/api/graph/generated"
	"cashflow/models"
	"context"
	"fmt"

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

	token, ok := ctx.Value(models.JwtKey).(string)
	if !ok || token == "" {
		return nil, fmt.Errorf("access denied")
	}

	user, err := (*d.AuthService).ValidateSession(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("access denied")
	}

	if user.Role != requires.String() {
		return nil, fmt.Errorf("access denied")
	}

	ctx = context.WithValue(ctx, models.UserId, user.ID)
	ctx = context.WithValue(ctx, models.UserRole, user.Role)
	return next(ctx)
}
