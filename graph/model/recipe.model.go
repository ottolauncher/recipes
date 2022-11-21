package model

import (
	pager "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID            primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name          string               `json:"name"`
	Slug          *string              `json:"slug,omitempty" bson:"slug,omitempty"`
	Timers        []string             `json:"timers"`
	Steps         []string             `json:"steps"`
	ImageURL      string               `json:"imageURL" bson:"imageURL"`
	OriginalURL   *string              `json:"originalURL" bson:"originalURL"`
	Ingredients   []*Ingredient        `json:"ingredients" bson:"ingredients"`
	IngredientIDs []primitive.ObjectID `json:"ingredient_ids,omitempty" bson:"ingredient_ids,omiempty"`
	Pagination    pager.PaginatedData  `json:"pagination,omitempty"`
}

func (r *Recipe) IsBaseModel() {}

func (r *Recipe) GetID() string {
	return r.ID.Hex()
}

func (r *Recipe) GetName() string {
	return r.Name
}

func (r *Recipe) GetSlug() *string {
	return r.Slug
}

func (r *Recipe) IsSearchRecipeResult() {}
