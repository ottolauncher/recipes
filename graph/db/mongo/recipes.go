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
}

func NewRecipeManager(d *mongo.Database) *RecipeManager {
	recipes := d.Collection("recipes")
	return &RecipeManager{Col: recipes}
}

func (tm *RecipeManager) Bulk(ctx context.Context, args []*model.NewRecipe) error {
	_, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	src := []interface{}{}

	for _, arg := range args {
		go func(v *model.NewRecipe) {
			slug := text.Slugify(v.Name)

			var ingredients []model.Ingredient

			for _, i := range v.Ingredients {
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
				"name":        v.Name,
				"slug":        &slug,
				"timers":      v.Timers,
				"steps":       v.Steps,
				"imageURL":    v.ImageURL,
				"originalURL": &v.OriginalURL,
				"ingredients": ingredients,
			}
			src = append(src, input)
		}(arg)
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
	opts := options.FindOptions{}
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

	search := bson.M{
		"$text": bson.M{
			"$search": query,
		},
	}

	var recipes []*model.Recipe
	cur, err := tm.Col.Find(l, search)

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
