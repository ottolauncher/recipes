package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ottolauncher/recipes/graph/model"
	"github.com/ottolauncher/recipes/preloads"
	"github.com/ottolauncher/recipes/utils/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IRecipe interface {
	Create(ctx context.Context, args model.NewRecipe) (*model.Recipe, error)
	Update(ctx context.Context, args model.UpdateRecipe) (*model.Recipe, error)
	Delete(ctx context.Context, filter map[string]interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error)
	All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error)
	Search(ctx context.Context, query string, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error)
}

type RecipeManager struct {
	Col *mongo.Collection
}

func NewRecipeManager(d *mongo.Database) *RecipeManager {
	recipes := d.Collection("recipes")
	return &RecipeManager{Col: recipes}
}

func (tm *RecipeManager) Create(ctx context.Context, args model.NewRecipe) (*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	Recipe := model.Recipe{
		Name:        args.Name,
		Slug:        &slug,
		Times:       args.Timers,
		Steps:       args.Steps,
		ImageURL:    &args.ImageURL,
		OriginalURL: &args.OriginalURL,
		Ingredients: ingredients,
	}
	res, err := tm.Col.InsertOne(l, Recipe)
	if err != nil {
		return nil, err
	}
	Recipe.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return &Recipe, nil
}

func (tm *RecipeManager) Update(ctx context.Context, args model.Recipe) (*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	Recipe := model.Recipe{
		Name:        args.Name,
		Note:        args.Name,
		Slug:        slug,
		Times:       args.Timers,
		Steps:       args.Steps,
		ImageURL:    args.ImageURL,
		OriginalURL: args.OriginalURL,
		Ingredients: args.Ingredients,
	}

	res, err := tm.Col.UpdateByID(l, args.ID, Recipe)
	if err != nil {
		return nil, err
	}
	Recipe.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	return &Recipe, nil
}

func (tm *RecipeManager) Delete(ctx context.Context, filter map[string]interface{}) error {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	if value, ok := filter["id"]; ok {
		pk, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return err
		}
		_, err = tm.Col.DeleteOne(l, pk)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tm *RecipeManager) Get(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error) {
	load := preloads.GetPreloads(ctx)
	projections := primitive.M{}

	for _, p := range load {
		projections[p] = 1
	}
	opts := options.FindOne().SetProjection(projections)
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var Recipe model.Recipe
	err := tm.Col.FindOne(l, filter, opts).Decode(&Recipe)
	if err != nil {
		return nil, err
	}
	return &Recipe, nil

}

func (tm *RecipeManager) All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 2000*time.Millisecond)
	defer cancel()
	load := preloads.GetPreloads(ctx)
	projections := primitive.M{}
	for _, p := range load {
		projections[p] = 1
	}
	opts := options.FindOptions{
		Projection: projections,
	}
	opts.SetLimit(int64(limit))

	var Recipes []*model.Recipe
	cur, err := tm.Col.Find(l, filter, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &Recipes); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return Recipes, nil
	}
	_ = cur.Close(l)
	if len(Recipes) == 0 {
		return Recipes, mongo.ErrNoDocuments
	}
	return Recipes, nil
}

func (tm *RecipeManager) Search(ctx context.Context, query string, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()
	load := preloads.GetPreloads(ctx)
	projections := primitive.M{}
	for _, p := range load {
		projections[p] = 1
	}
	opts := options.FindOptions{
		Projection: projections,
	}

	search := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	var Recipes []*model.Recipe
	cur, err := tm.Col.Find(l, search, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &Recipes); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return Recipes, nil
	}
	_ = cur.Close(l)
	if len(Recipes) == 0 {
		return Recipes, mongo.ErrNoDocuments
	}
	return Recipes, nil
}
