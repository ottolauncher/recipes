// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type BaseModel interface {
	IsBaseModel()
	GetID() string
	GetName() string
	GetSlug() *string
}

type SearchRecipeResult interface {
	IsSearchRecipeResult()
}

type Ingredient struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Slug     *string `json:"slug"`
	Type     string  `json:"type"`
	Quantity int     `json:"quantity"`
}

func (Ingredient) IsBaseModel()          {}
func (this Ingredient) GetID() string    { return this.ID }
func (this Ingredient) GetName() string  { return this.Name }
func (this Ingredient) GetSlug() *string { return this.Slug }

func (Ingredient) IsSearchRecipeResult() {}

type NewIngredient struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type NewRecipe struct {
	Name        string           `json:"name"`
	Timers      []string         `json:"timers"`
	Steps       []string         `json:"steps"`
	ImageURL    string           `json:"imageURL"`
	OriginalURL string           `json:"originalURL"`
	Ingredients []*NewIngredient `json:"ingredients"`
}

type Recipe struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Slug        *string       `json:"slug"`
	Timers      []string      `json:"timers"`
	Steps       []string      `json:"steps"`
	ImageURL    string        `json:"imageURL"`
	OriginalURL string        `json:"originalURL"`
	Ingredients []*Ingredient `json:"ingredients"`
}

func (Recipe) IsBaseModel()          {}
func (this Recipe) GetID() string    { return this.ID }
func (this Recipe) GetName() string  { return this.Name }
func (this Recipe) GetSlug() *string { return this.Slug }

func (Recipe) IsSearchRecipeResult() {}

type UpdateIngredient struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type UpdateRecipe struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Timers      []string            `json:"timers"`
	Steps       []string            `json:"steps"`
	ImageURL    string              `json:"imageURL"`
	OriginalURL string              `json:"originalURL"`
	Ingredients []*UpdateIngredient `json:"ingredients"`
}
