package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/ottolauncher/recipes/graph/generated"
	"github.com/ottolauncher/recipes/graph/model"
)

// CreateIngredient is the resolver for the createIngredient field.
func (r *mutationResolver) CreateIngredient(ctx context.Context, input model.NewIngredient) (*model.Ingredient, error) {
	panic(fmt.Errorf("not implemented: CreateIngredient - createIngredient"))
}

// UpdateIngredient is the resolver for the updateIngredient field.
func (r *mutationResolver) UpdateIngredient(ctx context.Context, input *model.UpdateIngredient) (*model.Ingredient, error) {
	panic(fmt.Errorf("not implemented: UpdateIngredient - updateIngredient"))
}

// DeleteIngredient is the resolver for the deleteIngredient field.
func (r *mutationResolver) DeleteIngredient(ctx context.Context, filter map[string]interface{}) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteIngredient - deleteIngredient"))
}

// CreateRecipe is the resolver for the createRecipe field.
func (r *mutationResolver) CreateRecipe(ctx context.Context, input model.NewRecipe) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented: CreateRecipe - createRecipe"))
}

// UpdateRecipe is the resolver for the updateRecipe field.
func (r *mutationResolver) UpdateRecipe(ctx context.Context, input model.UpdateRecipe) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented: UpdateRecipe - updateRecipe"))
}

// DeleteRecipe is the resolver for the deleteRecipe field.
func (r *mutationResolver) DeleteRecipe(ctx context.Context, filter map[string]interface{}) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteRecipe - deleteRecipe"))
}

// Ingredient is the resolver for the ingredient field.
func (r *queryResolver) Ingredient(ctx context.Context, filter map[string]interface{}) (*model.Ingredient, error) {
	panic(fmt.Errorf("not implemented: Ingredient - ingredient"))
}

// Ingredients is the resolver for the ingredients field.
func (r *queryResolver) Ingredients(ctx context.Context, filter map[string]interface{}, limit *int, page *int) ([]*model.Ingredient, error) {
	panic(fmt.Errorf("not implemented: Ingredients - ingredients"))
}

// Recipe is the resolver for the recipe field.
func (r *queryResolver) Recipe(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error) {
	panic(fmt.Errorf("not implemented: Recipe - recipe"))
}

// Recipes is the resolver for the recipes field.
func (r *queryResolver) Recipes(ctx context.Context, filter map[string]interface{}, limit *int, page *int) ([]*model.Recipe, error) {
	panic(fmt.Errorf("not implemented: Recipes - recipes"))
}

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, filter map[string]interface{}, limit *int, page *int) ([]model.SearchRecipeResult, error) {
	panic(fmt.Errorf("not implemented: Search - search"))
}

// Recipe is the resolver for the recipe field.
func (r *subscriptionResolver) Recipe(ctx context.Context) (<-chan *model.Recipe, error) {
	panic(fmt.Errorf("not implemented: Recipe - recipe"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
