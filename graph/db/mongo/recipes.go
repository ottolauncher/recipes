package db

import (
	"context"
	"fmt"
	"time"

	pager "github.com/gobeam/mongo-go-pagination"
	"github.com/ottolauncher/recipes/graph/model"
	"github.com/ottolauncher/recipes/utils/text"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRecipe interface {
	Create(ctx context.Context, args *model.NewRecipe) error
	Bulk(ctx context.Context, args []*model.NewRecipe) error
	Update(ctx context.Context, args *model.UpdateRecipe) error
	Delete(ctx context.Context, filter map[string]interface{}) error
	Get(ctx context.Context, filter map[string]interface{}) (*model.Recipe, error)
	All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error)
	Search(ctx context.Context, query string, limit int, page int) ([]*model.Recipe, error)
}

type RecipeManager struct {
	Col *mongo.Collection
	DB  *mongo.Database
}

func NewRecipeManager(d *mongo.Database) *RecipeManager {
	recipes := d.Collection("recipes")
	return &RecipeManager{Col: recipes, DB: d}
}

func (tm *RecipeManager) Bulk(ctx context.Context, args []*model.NewRecipe) error {
	_, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	src := []interface{}{}

	for _, v := range args {
		lsrc := []interface{}{}
		slug := text.Slugify(v.Name)
		id := primitive.NewObjectID()

		for _, i := range v.Ingredients {
			slg := text.Slugify(i.Name)
			lsrc = append(lsrc, bson.M{
				"name":      i.Name,
				"slug":      &slg,
				"type":      i.Type,
				"quantity":  i.Quantity,
				"recipe_id": id,
			})
		}

		res, err := tm.DB.Collection("ingredients").InsertMany(context.TODO(), lsrc)
		if err != nil {
			return err
		}

		input := bson.M{
			"_id":           id,
			"name":          v.Name,
			"slug":          &slug,
			"timers":        v.Timers,
			"steps":         v.Steps,
			"imageURL":      v.ImageURL,
			"originalURL":   &v.OriginalURL,
			"ingredientIDs": res.InsertedIDs,
		}
		src = append(src, input)
	}

	_, err := tm.Col.InsertMany(context.TODO(), src)
	if err != nil {
		return err
	}

	return nil
}

func (tm *RecipeManager) Create(ctx context.Context, args *model.NewRecipe) error {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)

	var ingredients []model.Ingredient

	for _, i := range args.Ingredients {
		slg := text.Slugify(i.Name)
		ingredients = append(ingredients, model.Ingredient{
			ID:       primitive.NewObjectID(),
			Name:     i.Name,
			Slug:     &slg,
			Type:     i.Type,
			Quantity: i.Quantity,
		})
	}
	input := bson.M{
		"name":        args.Name,
		"slug":        &slug,
		"timers":      args.Timers,
		"steps":       args.Steps,
		"imageURL":    args.ImageURL,
		"originalURL": &args.OriginalURL,
		"ingredients": ingredients,
	}

	_, err := tm.Col.InsertOne(l, input)
	if err != nil {
		return err
	}

	return nil
}

func (tm *RecipeManager) Update(ctx context.Context, args *model.UpdateRecipe) error {
	l, cancel := context.WithTimeout(ctx, 350*time.Millisecond)
	defer cancel()
	slug := text.Slugify(args.Name)
	var ingredients []model.Ingredient

	for _, i := range args.Ingredients {
		slg := text.Slugify(i.Name)
		ingredients = append(ingredients, model.Ingredient{
			Name:     i.Name,
			Slug:     &slg,
			Type:     i.Type,
			Quantity: i.Quantity,
		})
	}

	recipe := bson.D{
		{"$set",
			bson.D{{"name", args.Name},
				{"slug", slug},
				{"timers", args.Timers},
				{"steps", args.Steps},
				{"imageURL", args.ImageURL},
				{"originalURL", &args.OriginalURL},
				{"ingredients", ingredients},
			},
		},
	}

	id, err := primitive.ObjectIDFromHex(args.ID)
	if err != nil {
		return err
	}
	_, err = tm.Col.UpdateByID(l, id, recipe)
	if err != nil {
		return err
	}
	return nil
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
	l, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	var (
		recipe model.Recipe
		err    error
	)

	if id, ok := filter["id"]; ok {
		i, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", id))
		if err != nil {
			return nil, err
		}
		err = tm.Col.FindOne(l, bson.M{"_id": i}).Decode(&recipe)
		if err != nil {
			return nil, err
		}
	} else {

		err = tm.Col.FindOne(l, filter).Decode(&recipe)
		if err != nil {
			return nil, err
		}
	}

	return &recipe, nil

}

func (tm *RecipeManager) All(ctx context.Context, filter map[string]interface{}, limit int, page int) ([]*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	matchStage := bson.M{"$match": filter}
	lookupStage := bson.M{"$lookup": bson.M{"from": "ingredients", "localField": "_id", "foreignField": "recipe_id", "as": "ingredients"}}
	// matchStage := bson.D{{"$match"}}

	var recipes []*model.Recipe
	cur, err := pager.New(tm.Col).Context(l).Limit(int64(limit)).Page(int64(page)).Aggregate(lookupStage, matchStage)

	if err != nil {
		return nil, err
	}
	for _, raw := range cur.Data {
		var recipe *model.Recipe
		if marshallErr := bson.Unmarshal(raw, &recipe); marshallErr == nil {
			recipe.Pagination = *cur
			recipes = append(recipes, recipe)
		}
	}
	if len(recipes) == 0 {
		return recipes, mongo.ErrNoDocuments
	}

	return recipes, nil
}

func (tm *RecipeManager) Search(ctx context.Context, query string, limit int, page int) ([]*model.Recipe, error) {
	l, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	matchStage := bson.M{"$match": bson.M{"$text": bson.M{"$search": query}}}
	lookupStage := bson.M{"$lookup": bson.M{"from": "ingredients", "localField": "_id", "foreignField": "recipe_id", "as": "ingredients"}}

	var recipes []*model.Recipe
	cur, err := pager.New(tm.Col).Context(l).Limit(int64(limit)).Page(int64(page)).Aggregate(lookupStage, matchStage)

	if err != nil {
		return nil, err
	}
	for _, raw := range cur.Data {
		var recipe *model.Recipe
		if marshallErr := bson.Unmarshal(raw, &recipe); marshallErr == nil {
			recipe.Pagination = *cur
			recipes = append(recipes, recipe)
		}
	}

	if len(recipes) == 0 {
		return recipes, mongo.ErrNoDocuments
	}
	return recipes, nil
}
