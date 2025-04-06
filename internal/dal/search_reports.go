package dal

import (
	"fmt"
	"strings"

	"hot-coffee/internal/utils"
	"hot-coffee/models"

	"github.com/lib/pq"
)

type ReportRepository interface {
	SearchReports(query string, filters []string, minPrice, maxPrice float64) (*SearchResult, error)
}

type reportRepo struct {
	path string
}

func NewReportRepo(path string) *reportRepo {
	return &reportRepo{path: path}
}

type SearchResult struct {
	MenuItems []models.MenuItem          `json:"menu_items"`
	Orders    []models.OrderSearchResult `json:"orders"`
	Total     int                        `json:"total_matches"`
}

func (r *reportRepo) SearchReports(query string, filters []string, minPrice, maxPrice float64) (*SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	sqlQuery := `
		SELECT menu_item_id, name, description, price, 
           ts_rank_cd(to_tsvector('english', name || ' ' || description || ' ' || menu_item_id::text), plainto_tsquery($1)) as relevance
    FROM menu_items
    WHERE to_tsvector('english', name || ' ' || description || ' ' || menu_item_id::text) @@ plainto_tsquery($1)
	`
	args := []interface{}{query}

	if minPrice > 0 {
		sqlQuery += " AND price >= $2"
		args = append(args, minPrice)
	}
	if maxPrice > 0 {
		sqlQuery += " AND price <= $3"
		args = append(args, maxPrice)
	}
	sqlQuery += " ORDER BY relevance DESC;"

	rows, err := utils.DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results SearchResult
	for rows.Next() {
		var menuItem models.MenuItem
		var relevance float64
		err := rows.Scan(&menuItem.ID, &menuItem.Name, &menuItem.Description, &menuItem.Price, &relevance)
		if err != nil {
			return nil, err
		}
		menuItem.Relevance = relevance
		results.MenuItems = append(results.MenuItems, menuItem)
	}

	if contains(filters, "orders") || contains(filters, "all") {
		orderQuery := `
        SELECT o.order_id, o.customer_name, array_agg(oi.menu_item_id), o.total_amount, 
       ts_rank_cd(to_tsvector('english', o.customer_name), plainto_tsquery($1)) as relevance
FROM orders o
JOIN order_items oi ON o.order_id = oi.order_id
WHERE to_tsvector('english', o.customer_name) @@ plainto_tsquery($1)
OR oi.menu_item_id::text LIKE '%' || $1 || '%'
GROUP BY o.order_id
ORDER BY relevance DESC;

    `
		orderRows, err := utils.DB.Query(orderQuery, query)
		if err != nil {
			return nil, err
		}
		defer orderRows.Close()

		for orderRows.Next() {
			var order models.OrderSearchResult
			var items pq.StringArray
			err := orderRows.Scan(&order.ID, &order.CustomerName, &items, &order.Total, &order.Relevance)
			if err != nil {
				return nil, err
			}
			order.Items = items
			results.Orders = append(results.Orders, order)
		}
	}

	results.Total = len(results.MenuItems) + len(results.Orders)
	return &results, nil
}

func contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}
