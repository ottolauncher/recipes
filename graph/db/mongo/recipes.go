package db

import (
	"context"
	"fmt"
	"log"
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
	Create(ctx context.Context, args *model.NewRecipe) (*model.Recipe, error)
	Bulk(ctx context.Context, args []*model.NewRecipe) error
	Update(ctx context.Context, args *model.UpdateRecipe) (*model.Recipe, error)
	Delete(ctx context.Context, filter map[string]interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error)
	All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error)
	Search(ctx context.Context, query string, limit int, page int) ([]*model.Recipe, error)
}

type RecipeManager struct {
	Col *mongo.Collection
}

func NewRecipeManager(d *mongo.Database) *RecipeManager {
	recipes := d.Collection("recipes")
	return &RecipeManager{Col: recipes}
}

func (tm *RecipeManager) Bulk(ctx context.Context, args []*model.NewRecipe) error {
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	src := make([]interface{}, len(args))

	for _, arg := range args {
		slug := text.Slugify(arg.Name)

		var (
			timers      []*string
			steps       []*string
			ingredients []*model.Ingredient
		)

		for _, t := range arg.Timers {
			timers = append(timers, &t)
		}

		for _, s := range arg.Steps {
			steps = append(steps, &s)
		}

		for _, i := range arg.Ingredients {
			slg := text.Slugify(i.Name)
			ingredients = append(ingredients, &model.Ingredient{
				Name:     arg.Name,
				Slug:     &slg,
				Type:     i.Type,
				Quantity: i.Quantity,
			})
		}

		recipe := model.Recipe{
			Name:        arg.Name,
			Slug:        &slug,
			Timers:      timers,
			Steps:       steps,
			ImageURL:    arg.ImageURL,
			OriginalURL: &arg.OriginalURL,
			Ingredients: ingredients,
		}
		src = append(src, recipe)

	}

	_, err := tm.Col.InsertMany(l, src)
	if err != nil {
		return err
	}

	return nil

}

func (tm *RecipeManager) Create(ctx context.Context, args *model.NewRecipe) (*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	var (
		timers      []*string
		steps       []*string
		ingredients []*model.Ingredient
	)

	for _, t := range args.Timers {
		go func(u string) {
			timers = append(timers, &u)
		}(t)
	}

	for _, s := range args.Steps {
		go func(u string) {
			steps = append(steps, &u)
		}(s)
	}

	for _, i := range args.Ingredients {
		go func(v *model.NewIngredient) {
			slg := text.Slugify(v.Name)
			ingredients = append(ingredients, &model.Ingredient{
				ID:       primitive.NewObjectID().Hex(),
				Name:     v.Name,
				Slug:     &slg,
				Type:     v.Type,
				Quantity: v.Quantity,
			})
		}(i)
	}

	recipe := model.Recipe{
		Name:        args.Name,
		Slug:        &slug,
		Timers:      timers,
		Steps:       steps,
		ImageURL:    args.ImageURL,
		OriginalURL: &args.OriginalURL,
		Ingredients: ingredients,
	}
	log.Println(recipe)
	res, err := tm.Col.InsertOne(l, recipe)
	if err != nil {
		return nil, err
	}
	log.Println(res)
	recipe.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return &recipe, nil
}

func (tm *RecipeManager) Update(ctx context.Context, args *model.UpdateRecipe) (*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	var (
		timers      []*string
		steps       []*string
		ingredients []*model.Ingredient
	)

	for _, t := range args.Timers {
		timers = append(timers, &t)
	}

	for _, s := range args.Steps {
		steps = append(steps, &s)
	}

	for _, i := range args.Ingredients {
		slg := text.Slugify(i.Name)
		ingredients = append(ingredients, &model.Ingredient{
			Name:     args.Name,
			Slug:     &slg,
			Type:     i.Type,
			Quantity: i.Quantity,
		})
	}
	recipe := model.Recipe{
		Name:        args.Name,
		Slug:        &slug,
		Timers:      timers,
		Steps:       steps,
		ImageURL:    args.ImageURL,
		OriginalURL: &args.OriginalURL,
		Ingredients: ingredients,
	}

	res, err := tm.Col.UpdateByID(l, args.ID, recipe)
	if err != nil {
		return nil, err
	}
	recipe.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	return &recipe, nil
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

	var recipe model.Recipe
	err := tm.Col.FindOne(l, filter, opts).Decode(&recipe)
	if err != nil {
		return nil, err
	}
	return &recipe, nil

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

	var recipes []*model.Recipe
	cur, err := tm.Col.Find(l, filter, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &recipes); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return recipes, nil
	}
	_ = cur.Close(l)
	if len(recipes) == 0 {
		return recipes, mongo.ErrNoDocuments
	}
	return recipes, nil
}

func (tm *RecipeManager) Search(ctx context.Context, query string, limit int, page int) ([]*model.Recipe, error) {
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

	var recipes []*model.Recipe
	cur, err := tm.Col.Find(l, search, &opts)

	if err != nil {
		return nil, err
	}
	if err := cur.All(l, &recipes); err != nil {
		return nil, err
	}

	if err := cur.Err(); err != nil {
		return recipes, nil
	}
	_ = cur.Close(l)
	if len(recipes) == 0 {
		return recipes, mongo.ErrNoDocuments
	}
	return recipes, nil
}
