package models

type MenuItem struct {
	ID          string               `json:"menu_item_id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Ingredients []MenuItemIngredient `json:"ingredients"`
	Relevance   float64              `json:"relevance"`
}

type MenuItemIngredient struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}
