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

type Ingredient interface {
	Create(ctx context.Context, args model.NewIngredient) (*model.Ingredient, error)
	Update(ctx context.Context, args model.UpdateIngredient) (*model.Ingredient, error)
	Delete(ctx context.Context, filter map[string]interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) (*model.Ingredient, error)
	All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error)
	Search(ctx context.Context, query string, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error)
}

type IngredientManager struct {
	Col *mongo.Collection
}

func NewIngredientManager(d *mongo.Database) *IngredientManager {
	ingredients := d.Collection("ingredients")
	return &IngredientManager{Col: ingredients}
}

func (tm *IngredientManager) Create(ctx context.Context, args model.Ingredient) (*model.Ingredient, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	Ingredient := model.Ingredient{
		Name:     args.Name,
		Slug:     &slug,
		Type:     args.Type,
		Quantity: args.Quantity,
	}
	res, err := tm.Col.InsertOne(l, Ingredient)
	if err != nil {
		return nil, err
	}
	Ingredient.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return &Ingredient, nil
}

func (tm *IngredientManager) Update(ctx context.Context, args model.Ingredient) (*model.Ingredient, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	Ingredient := model.Ingredient{
		Name:     args.Name,
		Slug:     &slug,
		Type:     args.Type,
		Quantity: args.Quantity,
	}
	res, err := tm.Col.UpdateByID(l, args.ID, Ingredient)
	if err != nil {
		return nil, err
	}
	Ingredient.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	return &Ingredient, nil
}

func (tm *IngredientManager) Delete(ctx context.Context, filter map[string]interface{}) error {
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

func (tm *IngredientManager) Get(ctx context.Context, filter map[string]interface{}) (*model.Ingredient, error) {
	load := preloads.GetPreloads(ctx)
	projections := primitive.M{}

	for _, p := range load {
		projections[p] = 1
	}
	opts := options.FindOne().SetProjection(projections)
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var Ingredient model.Ingredient
	err := tm.Col.FindOne(l, filter, opts).Decode(&Ingredient)
	if err != nil {
		return nil, err
	}
	return &Ingredient, nil

}

func (tm *IngredientManager) All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error) {
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

	var Ingredients []*model.Ingredient
	cur, err := tm.Col.Find(l, filter, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &Ingredients); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return Ingredients, nil
	}
	_ = cur.Close(l)
	if len(Ingredients) == 0 {
		return Ingredients, mongo.ErrNoDocuments
	}
	return Ingredients, nil
}

func (tm *IngredientManager) Search(ctx context.Context, query string, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error) {
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

	var Ingredients []*model.Ingredient
	cur, err := tm.Col.Find(l, search, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &Ingredients); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return Ingredients, nil
	}
	_ = cur.Close(l)
	if len(Ingredients) == 0 {
		return Ingredients, mongo.ErrNoDocuments
	}
	return Ingredients, nil
}
