package model

type Ingredient struct {
	ID       string  `json:"id" bson:"_id"`
	Name     string  `json:"name"`
	Slug     *string `json:"slug,omitempty" bson:"slug,omitempty"`
	Type     string  `json:"type"`
	Quantity int     `json:"quantity"`
}

func (r *Ingredient) IssBaseModel() {}
