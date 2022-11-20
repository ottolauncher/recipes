package model

type Recipe struct {
	ID          string        `json:"id" bson:"_id"`
	Name        string        `json:"name"`
	Slug        *string       `json:"slug,omitempty" bson:"slug,omitempty"`
	Timers      *[]string     `json:"timers"`
	Steps       *[]string     `json:"steps"`
	ImageURL    *string       `json:"imageURL" bson:"imageURL"`
	OriginalURL *string       `json:"originalURL" bson:"originalURL"`
	Ingredients *[]Ingredient `json:"ingredients" bson:"ingredients"`
}

func (r *Recipe) IssBaseModel() {}
