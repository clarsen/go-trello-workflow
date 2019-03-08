//go:generate go run github.com/99designs/gqlgen
package handle_graphql_gqlgen

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string) (*User, error) {
	userClaims := ForContext(ctx)
	name := "a name"
	var email = "undefined"
	if userClaims != nil {
		email, _ = userClaims.String("email")
	}
	name = name + " " + email
	return &User{ID: id, Name: &name}, nil
}
