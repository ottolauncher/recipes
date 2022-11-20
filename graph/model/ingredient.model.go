package model

type Ingredient struct {
	ID       string  `json:"id" bson:"_id"`
	Name     string  `json:"name"`
	Slug     *string `json:"slug,omitempty" bson:"slug,omitempty"`
	Type     string  `json:"type"`
	Quantity string  `json:"quantity"`
}

func (i *Ingredient) IsBaseModel() {}

func (i *Ingredient) GetID() string {
	return i.ID
}

func (i *Ingredient) GetName() string {
	return i.Name
}

func (i *Ingredient) GetSlug() *string {
	return i.Slug
}

func (i *Ingredient) IsSearchRecipeResult() {}
