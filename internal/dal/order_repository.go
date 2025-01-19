package dal

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
	"os"
)

const pathOrder = "data/orders.json"

type OrderRepository interface {
	SaveOrder(models.Order) error
	SaveAll([]models.Order) error
	GetAll() ([]models.Order, error)
	Exists(orderID string) (bool, error)
}

type orderRepo struct {
	path string
}

func (r *orderRepo) SaveOrder(order models.Order) error {
	query := `INSERT INTO orders (customer_name, status, total_amount) VALUES ($1, $2, $3)`
	res, err := utils.DB.Exec(query, order.CustomerName, order.Status, order.TotalAmount)
	if err != nil {
		return err
	}
	fmt.Println("ff")
	fmt.Println(res.LastInsertId)
	return nil
}

func NewOrderRepo(path string) *orderRepo {
	return &orderRepo{path: path}
}

func (r *orderRepo) SaveAll(order []models.Order) error {
	jsonData, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		return err
	}
	_, err = os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	err = os.WriteFile(r.path, jsonData, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func (r *orderRepo) GetAll() ([]models.Order, error) {
	var result []models.Order
	_, err := os.Stat(r.path)
	if os.IsNotExist(err) {
		file, err := os.Create(r.path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}
	byteMenu, err := os.ReadFile(r.path)
	if len(byteMenu) == 0 {
		return result, err
	}

	if err := json.Unmarshal(byteMenu, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r *orderRepo) Exists(orderID string) (bool, error) {
	orderItems, err := r.GetAll()
	if err != nil {
		return false, err
	}
	for _, orderItem := range orderItems {
		if orderItem.ID == orderID {
			return true, nil
		}
	}
	return false, nil
}
