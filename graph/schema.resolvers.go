package graph

import (
	"context"
	"fmt"

	"github.com/dgryski/trifles/uuid"
	"github.com/ottolauncher/recipes/graph/generated"
	"github.com/ottolauncher/recipes/graph/model"
)

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

// ID is the resolver for the id field.
func (r *ingredientResolver) ID(ctx context.Context, obj *model.Ingredient) (string, error) {
	return obj.ID.Hex(), nil
}

// CreateIngredient is the resolver for the createIngredient field.
func (r *mutationResolver) CreateIngredient(ctx context.Context, input model.NewIngredient) (bool, error) {
	if err := r.IM.Create(ctx, &input); err != nil {
		return false, err
	}
	return true, nil
}

// BulkIngredient is the resolver for the bulkIngredient field.
func (r *mutationResolver) BulkIngredient(ctx context.Context, input []*model.NewIngredient) (bool, error) {
	if err := r.IM.Bulk(ctx, input); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateIngredient is the resolver for the updateIngredient field.
func (r *mutationResolver) UpdateIngredient(ctx context.Context, input *model.UpdateIngredient) (bool, error) {
	if err := r.IM.Update(ctx, input); err != nil {
		return false, err
	}
	return true, nil
}

// DeleteIngredient is the resolver for the deleteIngredient field.
func (r *mutationResolver) DeleteIngredient(ctx context.Context, filter map[string]interface{}) (bool, error) {
	if err := r.IM.Delete(ctx, filter); err != nil {
		return false, err
	}
	return true, nil
}

// CreateRecipe is the resolver for the createRecipe field.
func (r *mutationResolver) CreateRecipe(ctx context.Context, input model.NewRecipe) (bool, error) {
	if err := r.RM.Create(ctx, &input); err != nil {
		return false, err
	}
	return true, nil
}

// BulkRecipe is the resolver for the bulkRecipe field.
func (r *mutationResolver) BulkRecipe(ctx context.Context, input []*model.NewRecipe) (bool, error) {
	if err := r.RM.Bulk(ctx, input); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateRecipe is the resolver for the updateRecipe field.
func (r *mutationResolver) UpdateRecipe(ctx context.Context, input model.UpdateRecipe) (bool, error) {
	if err := r.RM.Update(ctx, &input); err != nil {
		return false, err
	}
	return true, nil
}

// DeleteRecipe is the resolver for the deleteRecipe field.
func (r *mutationResolver) DeleteRecipe(ctx context.Context, filter map[string]interface{}) (bool, error) {
	if err := r.RM.Delete(ctx, filter); err != nil {
		return false, err
	}
	return true, nil
}

// Ingredient is the resolver for the ingredient field.
func (r *queryResolver) Ingredient(ctx context.Context, filter map[string]interface{}) (*model.Ingredient, error) {
	res, err := r.IM.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Ingredients is the resolver for the ingredients field.
func (r *queryResolver) Ingredients(ctx context.Context, filter map[string]interface{}, limit *int, page *int) ([]*model.Ingredient, error) {
	res, err := r.IM.All(ctx, filter, *limit, *page)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Recipe is the resolver for the recipe field.
func (r *queryResolver) Recipe(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error) {
	res, err := r.RM.Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Recipes is the resolver for the recipes field.
func (r *queryResolver) Recipes(ctx context.Context, filter map[string]interface{}, limit *int, page *int) ([]*model.Recipe, error) {
	res, err := r.RM.All(ctx, filter, *limit, *page)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Search is the resolver for the search field.
func (r *queryResolver) Search(ctx context.Context, query string, limit *int, page *int) ([]model.SearchRecipeResult, error) {
	var (
		res         []model.SearchRecipeResult
		searchError []error
	)

	recipes, rerr := r.RM.Search(ctx, query, *limit, *page)
	if rerr != nil {
		searchError = append(searchError, rerr)
	}

	ingredients, ierr := r.IM.Search(ctx, query, *limit, *page)
	if ierr != nil {
		searchError = append(searchError, ierr)
	}
	if len(searchError) > 0 {
		return nil, fmt.Errorf("%s", searchError)
	}
	if len(recipes) > 0 {
		for _, rcp := range recipes {
			res = append(res, rcp)
		}
	}

	if len(ingredients) > 0 {
		for _, i := range ingredients {
			res = append(res, i)
		}
	}
	return res, nil
}

// ID is the resolver for the id field.
func (r *recipeResolver) ID(ctx context.Context, obj *model.Recipe) (string, error) {
	return obj.ID.Hex(), nil
}

// Recipe is the resolver for the recipe field.
func (r *subscriptionResolver) Recipe(ctx context.Context) (<-chan []*model.Recipe, error) {
	id := uuid.UUIDv4()
	recipes := make(chan []*model.Recipe, 1)

	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.RecipeObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()

	r.RecipeObservers[id] = recipes
	r.mu.Unlock()
	r.RecipeObservers[id] <- r.Recipes
	return recipes, nil
}

// Ingredient returns generated.IngredientResolver implementation.
func (r *Resolver) Ingredient() generated.IngredientResolver { return &ingredientResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Recipe returns generated.RecipeResolver implementation.
func (r *Resolver) Recipe() generated.RecipeResolver { return &recipeResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type ingredientResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type recipeResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *ingredientResolver) Quantity(ctx context.Context, obj *model.Ingredient) (string, error) {
	return obj.Quantity, nil
}
