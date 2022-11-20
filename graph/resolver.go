//go:generate go run github.com/99designs/gqlgen generate
package graph

import db "github.com/ottolauncher/recipes/graph/db/mongo"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RM *db.RecipeManager
	IM *db.IngredientManager
}
