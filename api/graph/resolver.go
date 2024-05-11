package graph

import "cashflow/models"

// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	*models.App
}
