package dal

import (
	"encoding/json"

	"os"

	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type InventoryRepository interface {
	SaveAll(item []models.InventoryItem) error
	GetAll() ([]models.InventoryItem, error)
	Exists(id string) (bool, error)
	AddItem(item models.InventoryItem) error
	DeleteItem(id string) error
	UpdateItem(item models.InventoryItem) error
	CheckInventory(items []models.OrderItem) (bool, error)
	UpdateInventory(items []models.OrderItem) error
	GetLeftovers(sortBy string, offset, limit int) ([]models.InventoryItem, int, error)
}

type inventoryRepo struct {
	path string
}

func NewInventoryRepo(path string) *inventoryRepo {
	return &inventoryRepo{path: path}
}

func (r *inventoryRepo) AddItem(item models.InventoryItem) error {
	_, err := utils.DB.Exec(`INSERT INTO inventory(ingredient_id, name, quantity, unit, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		item.IngredientID, item.Name, item.Quantity, item.Unit, item.CreatedAt, item.UpdatedAt)
	return err
}

func (r *inventoryRepo) DeleteItem(id string) error {
	query := `DELETE FROM inventory WHERE ingredient_id = $1`
	_, err := utils.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *inventoryRepo) SaveAll(item []models.InventoryItem) error {
	jsonData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return err
		}
		file.Close()
	}
	err = os.WriteFile(r.path, jsonData, 0o644)
	if err != nil {
		return err
	}

	return nil
}

func (r *inventoryRepo) GetAll() ([]models.InventoryItem, error) {
	query := `SELECT ingredient_id, name, quantity, unit, created_at, updated_at FROM inventory;`

	rows, err := utils.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var inventory models.InventoryItem
	var inventoryData []models.InventoryItem
	for rows.Next() {
		err := rows.Scan(&inventory.IngredientID, &inventory.Name,
			&inventory.Quantity, &inventory.Unit,
			&inventory.CreatedAt, &inventory.UpdatedAt)
		if err != nil {
			return nil, err
		}
		inventoryData = append(inventoryData, inventory)
	}

	return inventoryData, nil
}

func (r *inventoryRepo) Exists(id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM inventory WHERE ingredient_id = $1)`
	err := utils.DB.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *inventoryRepo) UpdateItem(item models.InventoryItem) error {
	var quantity float64
	err := utils.DB.QueryRow(`SELECT quantity FROM inventory WHERE ingredient_id = $1`, item.IngredientID).Scan(&quantity)
	if err != nil {
		return err
	}
	if quantity != item.Quantity {
		query := `INSERT INTO inventory_transactions (ingredient_id, old_quantity, new_quantity, unit) VALUES ($1, $2, $3, $4)`
		_, err = utils.DB.Exec(query, item.IngredientID, quantity, item.Quantity, item.Unit)
		if err != nil {
			return err
		}
	}

	query := `UPDATE inventory SET name = $1, quantity = $2, unit = $3, updated_at = $4 WHERE ingredient_id = $5`
	_, err = utils.DB.Exec(query, item.Name, item.Quantity, item.Unit, item.UpdatedAt, item.IngredientID)
	return err
}

func (r *inventoryRepo) CheckInventory(items []models.OrderItem) (bool, error) {
	for _, item := range items {
		query := `
			SELECT mi.ingredient_id, mi.quantity, i.quantity 
			FROM menu_item_ingredients mi
			JOIN inventory i ON mi.ingredient_id = i.ingredient_id
			WHERE mi.menu_item_id = $1
		`
		rows, err := utils.DB.Query(query, item.MenuItemID)
		if err != nil {
			return false, err
		}
		defer rows.Close()

		for rows.Next() {
			var ingredientID string
			var requiredQuantity, availableQuantity float64
			err = rows.Scan(&ingredientID, &requiredQuantity, &availableQuantity)
			if err != nil {
				return false, err
			}

			totalRequired := requiredQuantity * float64(item.Quantity)
			if totalRequired > availableQuantity {
				return false, nil
			}
		}
	}

	return true, nil
}

func (r *inventoryRepo) UpdateInventory(items []models.OrderItem) error {
	for _, item := range items {
		query := `
		SELECT ingredient_id, quantity
		FROM menu_item_ingredients
		WHERE menu_item_id = $1;
		`
		rows, err := utils.DB.Query(query, item.MenuItemID)
		if err != nil {
			return err
		}
		defer rows.Close()

		updateQuery := `
		UPDATE inventory 
		SET quantity = quantity - $1
		WHERE ingredient_id = $2;
		`

		for rows.Next() {
			var ingredientID string
			var requiredQuantity float64

			if err := rows.Scan(&ingredientID, &requiredQuantity); err != nil {
				return err
			}

			totalQuantity := item.Quantity * int(requiredQuantity)

			_, err := utils.DB.Exec(updateQuery, totalQuantity, ingredientID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *inventoryRepo) GetLeftovers(sortBy string, offset, limit int) ([]models.InventoryItem, int, error) {
	query := `SELECT name, quantity FROM inventory`
	countQuery := `SELECT COUNT(*) FROM inventory`

	switch sortBy {
	case "quantity":
		query += " ORDER BY quantity DESC"
	default:
		query += " ORDER BY ingredient_id"
	}

	query += " OFFSET $1 LIMIT $2"

	rows, err := utils.DB.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []models.InventoryItem
	for rows.Next() {
		var item models.InventoryItem
		if err := rows.Scan(&item.Name, &item.Quantity); err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	var total int
	err = utils.DB.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
