package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ottolauncher/recipes/graph/model"
	"github.com/ottolauncher/recipes/utils/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ingredient interface {
	Create(ctx context.Context, args *model.NewIngredient) error
	Bulk(ctx context.Context, args []*model.NewIngredient) error
	Update(ctx context.Context, args *model.UpdateIngredient) error
	Delete(ctx context.Context, filter map[string]interface{}) error

	Get(ctx context.Context, filter map[string]interface{}) (*model.Ingredient, error)
	All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error)
	Search(ctx context.Context, query string, limit int, page int) ([]*model.Ingredient, error)
}

type IngredientManager struct {
	Col *mongo.Collection
}

func NewIngredientManager(d *mongo.Database) *IngredientManager {
	ingredients := d.Collection("ingredients")
	return &IngredientManager{Col: ingredients}
}

func (im *IngredientManager) Bulk(ctx context.Context, args []*model.NewIngredient) error {
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	src := []interface{}{}

	for _, args := range args {
		slug := text.Slugify(args.Name)
		ingredient := model.Ingredient{
			Name:     args.Name,
			Slug:     &slug,
			Type:     args.Type,
			Quantity: args.Quantity,
		}

		src = append(src, ingredient)

	}

	_, err := im.Col.InsertMany(l, src)
	if err != nil {
		return err
	}

	return nil
}
func (tm *IngredientManager) Create(ctx context.Context, args *model.NewIngredient) error {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	ingredient := model.Ingredient{
		Name:     args.Name,
		Slug:     &slug,
		Type:     args.Type,
		Quantity: args.Quantity,
	}
	_, err := tm.Col.InsertOne(l, ingredient)
	if err != nil {
		return err
	}
	return nil
}

func (tm *IngredientManager) Update(ctx context.Context, args *model.UpdateIngredient) error {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	ingredient := model.Ingredient{
		Name:     args.Name,
		Slug:     &slug,
		Type:     args.Type,
		Quantity: args.Quantity,
	}
	_, err := tm.Col.UpdateByID(l, args.ID, ingredient)
	if err != nil {
		return err
	}

	return nil
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
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var (
		ingredient model.Ingredient
		err        error
	)
	if id, ok := filter["id"]; ok {
		i, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", id))
		if err != nil {
			return nil, err
		}
		err = tm.Col.FindOne(l, bson.M{"_id": i}).Decode(&ingredient)
		if err != nil {
			return nil, err
		}
	} else {

		err = tm.Col.FindOne(l, filter).Decode(&ingredient)
		if err != nil {
			return nil, err
		}
	}
	return &ingredient, nil

}

func (tm *IngredientManager) All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Ingredient, error) {
	l, cancel := context.WithTimeout(ctx, 2000*time.Millisecond)
	defer cancel()

	opts := options.FindOptions{}
	opts.SetLimit(int64(limit))

	var ingredients []*model.Ingredient
	cur, err := tm.Col.Find(l, filter, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &ingredients); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return ingredients, nil
	}
	_ = cur.Close(l)
	if len(ingredients) == 0 {
		return ingredients, mongo.ErrNoDocuments
	}
	return ingredients, nil
}

func (tm *IngredientManager) Search(ctx context.Context, query string, limit int, page int) ([]*model.Ingredient, error) {
	l, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	search := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	var ingredients []*model.Ingredient
	cur, err := tm.Col.Find(l, search)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &ingredients); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return ingredients, nil
	}
	_ = cur.Close(l)
	if len(ingredients) == 0 {
		return ingredients, mongo.ErrNoDocuments
	}
	return ingredients, nil
}
