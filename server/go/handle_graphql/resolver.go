//go:generate go run github.com/99designs/gqlgen
package handle_graphql

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Tasks(ctx context.Context, dueBefore *int) ([]Task, error) {
	return []Task{
		Task{
			ID:    "an id",
			Title: "A title",
		},
	}, nil

}
