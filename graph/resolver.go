//go:generate go run github.com/99designs/gqlgen generate
package graph

import (
	"sync"

	db "github.com/ottolauncher/recipes/graph/db/mongo"
	"github.com/ottolauncher/recipes/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// TODO use redis subscription instead of inmemory one
type Resolver struct {
	RM              *db.RecipeManager
	IM              *db.IngredientManager
	Recipes         []*model.Recipe
	RecipeObservers map[string]chan []*model.Recipe
	mu              sync.Mutex
}
