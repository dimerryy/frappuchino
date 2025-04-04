package dal

import (
	"database/sql"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type MenuRepository interface {
	DeleteMenuItem(menuItemID string) error
	GetAll() ([]models.MenuItem, error)
	Exists(menuID string) (bool, error)
	GetMenuItemPrice(menuItemID string) (float64, error)
	SaveMenuItem(menuItem models.MenuItem) error
	Update(menu models.MenuItem) error
}

type menuRepo struct {
	path string
}

func NewMenuRepo(path string) *menuRepo {
	return &menuRepo{path: path}
}

func (r *menuRepo) DeleteMenuItem(menuItemID string) error {
	tx, err := utils.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM menu_item_ingredients WHERE menu_item_id = $1`, menuItemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM menu_items WHERE menu_item_id = $1`, menuItemID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *menuRepo) GetAll() ([]models.MenuItem, error) {
	var menuItems []models.MenuItem
	menuItemMap := make(map[string]*models.MenuItem)

	query := `
	SELECT 
		m.menu_item_id, m.name, m.description, m.price, 
		mi.ingredient_id, mi.quantity
	FROM menu_items m
	LEFT JOIN menu_item_ingredients mi ON m.menu_item_id = mi.menu_item_id;
	`

	rows, err := utils.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menuID, name, description string
		var price float64
		var ingredientID sql.NullString
		var quantity sql.NullFloat64

		err := rows.Scan(&menuID, &name, &description, &price, &ingredientID, &quantity)
		if err != nil {
			return nil, err
		}
		menuItem, exists := menuItemMap[menuID]
		if !exists {
			menuItem = &models.MenuItem{
				ID:          menuID,
				Name:        name,
				Description: description,
				Price:       price,
				Ingredients: []models.MenuItemIngredient{},
			}
			menuItemMap[menuID] = menuItem
		}
		if ingredientID.Valid {
			menuItem.Ingredients = append(menuItem.Ingredients, models.MenuItemIngredient{
				IngredientID: ingredientID.String,
				Quantity:     quantity.Float64,
			})
		}
	}
	for _, item := range menuItemMap {
		menuItems = append(menuItems, *item)
	}

	return menuItems, nil
}

func (r *menuRepo) Exists(menuID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM menu_items WHERE menu_item_id = $1)`
	err := utils.DB.QueryRow(query, menuID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *menuRepo) GetMenuItemPrice(menuItemID string) (float64, error) {
	var price float64
	err := utils.DB.QueryRow(`SELECT price FROM menu_items WHERE menu_item_id = $1`, menuItemID).Scan(&price)
	return price, err
}

func (r *menuRepo) SaveMenuItem(menuItem models.MenuItem) error {
	query := `INSERT INTO menu_items(menu_item_id, name, description, price) VALUES ($1, $2, $3, $4)`
	_, err := utils.DB.Exec(query, menuItem.ID, menuItem.Name, menuItem.Description, menuItem.Price)
	if err != nil {
		return err
	}
	for _, ingredient := range menuItem.Ingredients {
		ingredientQuery := `INSERT INTO menu_item_ingredients(menu_item_id, ingredient_id, quantity) VALUES ($1, $2, $3)`
		_, err = utils.DB.Exec(ingredientQuery, menuItem.ID, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *menuRepo) Update(menu models.MenuItem) error {
	tx, err := utils.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE menu_items 
		SET name = $1, description = $2, price = $3 
		WHERE menu_item_id = $4
	`
	_, err = tx.Exec(query, menu.Name, menu.Description, menu.Price, menu.ID)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM menu_item_ingredients WHERE menu_item_id = $1`
	_, err = tx.Exec(deleteQuery, menu.ID)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO menu_item_ingredients (menu_item_id, ingredient_id, quantity)
		VALUES ($1, $2, $3)
	`
	for _, ingredient := range menu.Ingredients {
		_, err := tx.Exec(insertQuery, menu.ID, ingredient.IngredientID, ingredient.Quantity)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
