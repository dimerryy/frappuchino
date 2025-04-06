package models

import "encoding/json"

type Order struct {
	ID               int         `json:"order_id"`
	CustomerName     string      `json:"customer_name"`
	Items            []OrderItem `json:"items"`
	Status           string      `json:"status"`
	CreatedAt        string      `json:"created_at"`
	TotalAmount      float64     `json:"total_amount"`
	UpdatedAt        string      `json:"updated_at"`
	LastStatusChange string      `json:"last_status_change"`
}

type OrderItem struct {
	MenuItemID    string          `json:"menu_item_id"`
	Quantity      int             `json:"quantity"`
	Price         float64         `json:"-"`
	Customization json.RawMessage `json:"customization,omitempty"`
}

type TotalSales struct {
	Sales float64 `json:"total_sales: "`
}

type OrderStatusHistory struct {
	ID        int    `json:"id"`
	OrderID   int    `json:"order_id"`
	Notes     string `json:"notes"`
	UpdatedAt string `json:"updated_at"`
}

type OrderSearchResult struct {
	ID           int      `json:"id"`
	CustomerName string   `json:"customer_name"`
	Total        float64  `json:"total_amount"`
	Items        []string `json:"items"`
	Relevance    float64  `json:"relevance"`
}

type BatchOrderRequest struct {
	Orders []Order `json:"orders"`
}

type BatchOrderResponse struct {
	ProcessedOrders []ProcessedOrder `json:"processed_orders"`
	Summary         BatchSummary     `json:"summary"`
}

type ProcessedOrder struct {
	OrderID      int     `json:"order_id,omitempty"`
	CustomerName string  `json:"customer_name"`
	Status       string  `json:"status"` // accepted or rejected
	Total        float64 `json:"total,omitempty"`
	Reason       string  `json:"reason,omitempty"`
}

type BatchSummary struct {
	TotalOrders      int              `json:"total_orders"`
	Accepted         int              `json:"accepted"`
	Rejected         int              `json:"rejected"`
	TotalRevenue     float64          `json:"total_revenue"`
	InventoryUpdates []InventoryUsage `json:"inventory_updates,omitempty"`
}

type InventoryUsage struct {
	IngredientID int    `json:"ingredient_id"`
	Name         string `json:"name"`
	QuantityUsed int    `json:"quantity_used"`
	Remaining    int    `json:"remaining"`
}
