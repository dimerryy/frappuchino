package dal

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type OrderRepository interface {
	SaveOrder(models.Order) error
	GetAll() ([]models.Order, error)
	OrderExists(orderID int) (bool, error)
	UpdateOrder(order models.Order) error
	DeleteOrder(orderID int) error
	CloseOrder(id int) error
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	GetOrdersGroupedByDay(month string) (map[string]interface{}, error)
	GetOrdersGroupedByMonth(year string) (map[string]interface{}, error)
}

type orderRepo struct {
	path string
}

func (r *orderRepo) SaveOrder(order models.Order) error {
	query := `INSERT INTO orders (customer_name, status,order_date,last_status_change, total_amount, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING order_id`
	var orderID int
	err := utils.DB.QueryRow(query, order.CustomerName, order.Status, order.CreatedAt, order.CreatedAt, order.TotalAmount, order.UpdatedAt).Scan(&orderID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query := `INSERT INTO order_items (order_id, menu_item_id, quantity, price, customization) 
				  VALUES ($1, $2, $3, $4, $5)`

		_, err := utils.DB.Exec(query, orderID, item.MenuItemID, item.Quantity, item.Price, string(item.Customization))
		if err != nil {
			return err
		}
	}

	return nil
}

func NewOrderRepo(path string) *orderRepo {
	return &orderRepo{path: path}
}

func (r *orderRepo) GetAll() ([]models.Order, error) {
	query := `
	SELECT 
		o.order_id, o.customer_name, o.status, o.order_date, 
		o.last_status_change, o.total_amount, o.updated_at,
		oi.menu_item_id, oi.quantity, oi.price, oi.customization
	FROM orders o
	LEFT JOIN order_items oi ON o.order_id = oi.order_id
	ORDER BY o.order_id;
	`
	rows, err := utils.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordersMap := make(map[int]*models.Order)

	for rows.Next() {
		var order models.Order
		var orderItem models.OrderItem
		var customizationJSON []byte
		var menuItemID sql.NullString
		var quantity sql.NullInt64
		var price sql.NullFloat64

		err := rows.Scan(
			&order.ID, &order.CustomerName, &order.Status, &order.CreatedAt,
			&order.LastStatusChange, &order.TotalAmount, &order.UpdatedAt,
			&menuItemID, &quantity, &price, &customizationJSON,
		)
		if err != nil {
			return nil, err
		}
		orderItem.MenuItemID = menuItemID.String
		orderItem.Quantity = int(quantity.Int64)
		orderItem.Price = price.Float64

		if len(customizationJSON) > 0 {
			if err := json.Unmarshal(customizationJSON, &orderItem.Customization); err != nil {
				return nil, fmt.Errorf("error unmarshaling customization: %w", err)
			}
		}

		if existingOrder, exists := ordersMap[order.ID]; exists {
			existingOrder.Items = append(existingOrder.Items, orderItem)
		} else {
			order.Items = []models.OrderItem{orderItem}
			ordersMap[order.ID] = &order
		}
	}

	var orders []models.Order
	for _, order := range ordersMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (r *orderRepo) OrderExists(orderID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM orders WHERE order_id = $1)`
	err := utils.DB.QueryRow(query, orderID).Scan(&exists)
	return exists, err
}

func (r *orderRepo) UpdateOrder(order models.Order) error {
	tx, err := utils.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE orders 
		SET customer_name = $1, status = $2, total_amount = $3, updated_at = $4
		WHERE order_id = $5
	`
	_, err = tx.Exec(query, order.CustomerName, order.Status, order.TotalAmount, order.UpdatedAt, order.ID)
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM order_items WHERE order_id = $1`
	_, err = tx.Exec(deleteQuery, order.ID)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO order_items (order_id, menu_item_id, quantity, price, customization)
		VALUES ($1, $2, $3, $4, $5::jsonb)
	`
	for _, item := range order.Items {
		_, err := tx.Exec(insertQuery, order.ID, item.MenuItemID, item.Quantity, item.Price, item.Customization)
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

func (r *orderRepo) CloseOrder(id int) error {
	_, err := utils.DB.Exec(`INSERT INTO order_status_history (order_id, old_status, new_status) VALUES ($1, $2, $3)`, id, "active", "closed")
	if err != nil {
		return err
	}

	_, err = utils.DB.Exec(`UPDATE orders SET status = 'closed' WHERE order_id = $1`, id)
	return err
}

func (r *orderRepo) DeleteOrder(orderID int) error {
	var status string

	err := utils.DB.QueryRow(`SELECT status FROM orders WHERE order_id = $1`, orderID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("order not found")
		}
		return err
	}

	if status == "closed" {
		return errors.New("cannot delete a closed order")
	}

	tx, err := utils.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM orders WHERE order_id = $1`, orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *orderRepo) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	query := `
		SELECT mi.name, SUM(oi.quantity) 
		FROM order_items oi
		JOIN menu_items mi ON oi.menu_item_id = mi.menu_item_id
		JOIN orders o ON oi.order_id = o.order_id
	`

	var args []interface{}
	if startDate != "" && endDate != "" {
		query += " WHERE o.order_date BETWEEN $1 AND $2"
		args = append(args, startDate, endDate)
	}

	query += " GROUP BY mi.name"

	rows, err := utils.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make(map[string]int)
	for rows.Next() {
		var name string
		var quantity int
		err := rows.Scan(&name, &quantity)
		if err != nil {
			return nil, err
		}
		items[name] = quantity
	}

	return items, nil
}

func (r *orderRepo) GetOrdersGroupedByDay(month string) (map[string]interface{}, error) {
	query := `
        SELECT 
    		EXTRACT(DAY FROM order_date)::int AS day, 
    	COUNT(*) 
		FROM orders 
		WHERE TO_CHAR(order_date, 'FMMonth') ILIKE $1 
		GROUP BY day 
		ORDER BY day;`

	rows, err := utils.DB.Query(query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]interface{}{
		"period":       "day",
		"month":        strings.ToLower(month),
		"orderedItems": []map[string]int{},
	}

	for rows.Next() {
		var day int
		var count int
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}
		result["orderedItems"] = append(result["orderedItems"].([]map[string]int), map[string]int{fmt.Sprintf("%d", day): count})
	}

	return result, nil
}

func (r *orderRepo) GetOrdersGroupedByMonth(year string) (map[string]interface{}, error) {
	query := `
        SELECT 
    		EXTRACT(MONTH FROM order_date)::int AS month_num,
    		TO_CHAR(order_date, 'Month') AS month_name,
    	COUNT(*) 
		FROM orders 
		WHERE EXTRACT(YEAR FROM order_date)::text = $1
		GROUP BY month_num, month_name
		ORDER BY month_num;
`

	rows, err := utils.DB.Query(query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]interface{}{
		"period":       "month",
		"year":         year,
		"orderedItems": []map[string]int{},
	}

	for rows.Next() {
		var monthNum int
		var month string
		var count int
		if err := rows.Scan(&monthNum, &month, &count); err != nil {
			return nil, err
		}
		result["orderedItems"] = append(result["orderedItems"].([]map[string]int), map[string]int{strings.TrimSpace(strings.ToLower(month)): count})
	}

	return result, nil
}
