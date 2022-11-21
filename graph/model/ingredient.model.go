package model

import (
	pager "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ingredient struct {
	ID         primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name       string              `json:"name"`
	Slug       *string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Type       string              `json:"type"`
	Quantity   string              `json:"quantity"`
	RecipeID   primitive.ObjectID  `json:"recipe_id" bson:"recipe_id,omitempty"`
	Pagination pager.PaginatedData `json:"pagination,omitempty"`
}

func (i *Ingredient) IsBaseModel() {}

func (i *Ingredient) GetID() string {
	return i.ID.Hex()
}

func (i *Ingredient) GetName() string {
	return i.Name
}

func (i *Ingredient) GetSlug() *string {
	return i.Slug
}

func (i *Ingredient) IsSearchRecipeResult() {}
