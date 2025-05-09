package service

import (
	"errors"

	"hot-coffee/internal/dal"
	"hot-coffee/internal/utils"
	"hot-coffee/models"
)

type OrderService interface {
	GetOrderItemById(id int) (models.Order, error)
	GetOrderItem() ([]models.Order, error)
	PostOrUpdate(order models.Order, id int) error
	UpdateOrderStatus(orderId int) error
	DeleteOrder(orderID int) error
	GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error)
	GetOrdersGroupedByDay(month string) (map[string]interface{}, error)
	GetOrdersGroupedByMonth(year string) (map[string]interface{}, error)
	ProcessBatchOrders(orders []models.Order) (*models.BatchOrderResponse, error)
}

type orderService struct {
	orderRepo     dal.OrderRepository
	menuRepo      dal.MenuRepository
	inventoryRepo dal.InventoryRepository
}

func NewOrderService(orderRepo dal.OrderRepository, menuRepo dal.MenuRepository, inventoryRepo dal.InventoryRepository) *orderService {
	return &orderService{orderRepo: orderRepo, menuRepo: menuRepo, inventoryRepo: inventoryRepo}
}

func (s *orderService) GetOrderItemById(id int) (models.Order, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return models.Order{}, err
	}
	for _, orderItem := range orderItems {
		if orderItem.ID == id {
			return orderItem, nil
		}
	}
	return models.Order{}, errors.New("inventory item not found")
}

func (s *orderService) GetOrderItem() ([]models.Order, error) {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return []models.Order{}, err
	}

	return orderItems, nil
}

func (s *orderService) UpdateOrderStatus(id int) error {
	orderItems, err := s.orderRepo.GetAll()
	if err != nil {
		return err
	}

	for i := range orderItems {
		if orderItems[i].ID == id {
			if orderItems[i].Status == "closed" {
				return errors.New("order is already closed")
			}

			if err = s.inventoryRepo.UpdateInventory(orderItems[i].Items); err != nil {
				return err
			}

			if err = s.orderRepo.CloseOrder(id); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("order item not found")
}

func (s *orderService) DeleteOrder(orderID int) error {
	return s.orderRepo.DeleteOrder(orderID)
}

func (s *orderService) PostOrUpdate(order models.Order, id int) error {
	order.ID = id
	if !IsOrderValid(order) {
		return errors.New("order is invalid")
	}
	sufficient, err := s.inventoryRepo.CheckInventory(order.Items)
	if err != nil {
		return err
	}
	if !sufficient {
		return errors.New("not enough inventory for order")
	}

	err = IsValidOrder(order, s.menuRepo, s.inventoryRepo)
	if err != nil {
		return err
	}

	var totalAmount float64
	for i := range order.Items {
		price, err := s.menuRepo.GetMenuItemPrice(order.Items[i].MenuItemID)
		if err != nil {
			return err
		}
		order.Items[i].Price = price
		totalAmount += price * float64(order.Items[i].Quantity)
	}

	now := getFormattedTime()

	if order.ID == 0 {
		order.LastStatusChange = now
		order.CreatedAt = now
		order.UpdatedAt = now
		order.TotalAmount = totalAmount
		order.Status = "active"
		_, err = s.orderRepo.SaveOrder(order)
		return err
	} else {
		exists, err := s.orderRepo.OrderExists(order.ID)
		if err != nil {
			return err
		}
		if !exists {
			return errors.New("order does not exist")
		}
		order.Status = "active"
		order.UpdatedAt = now
		order.TotalAmount = totalAmount
		err = s.orderRepo.UpdateOrder(order)
		return err
	}
}

func (s *orderService) GetNumberOfOrderedItems(startDate, endDate string) (map[string]int, error) {
	return s.orderRepo.GetNumberOfOrderedItems(startDate, endDate)
}

func (s *orderService) GetOrdersGroupedByDay(month string) (map[string]interface{}, error) {
	return s.orderRepo.GetOrdersGroupedByDay(month)
}

func (s *orderService) GetOrdersGroupedByMonth(year string) (map[string]interface{}, error) {
	return s.orderRepo.GetOrdersGroupedByMonth(year)
}

func (s *orderService) ProcessBatchOrders(orders []models.Order) (*models.BatchOrderResponse, error) {
	var response models.BatchOrderResponse
	tx, err := utils.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	for _, order := range orders {
		ok, err := s.inventoryRepo.CheckInventory(order.Items)
		if err != nil {
			return nil, err
		}

		if !ok {
			response.ProcessedOrders = append(response.ProcessedOrders, models.ProcessedOrder{
				OrderID:      0,
				CustomerName: order.CustomerName,
				Status:       "rejected",
				Reason:       "insufficient_inventory",
			})
			response.Summary.Rejected++
			continue
		}

		err = s.PostOrUpdate(order, 0)
		if err != nil {
			return nil, err
		}

		var orderID int
		var total float64
		err = utils.DB.QueryRow(`SELECT order_id, total_amount
									FROM orders WHERE customer_name = $1;`, order.CustomerName).Scan(&orderID, &total)
		if err != nil {
			return nil, err
		}
		response.ProcessedOrders = append(response.ProcessedOrders, models.ProcessedOrder{
			OrderID:      orderID,
			CustomerName: order.CustomerName,
			Status:       "accepted",
			Total:        total,
		})
		response.Summary.Accepted++
		response.Summary.TotalRevenue += total
	}

	response.Summary.TotalOrders = len(orders)
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &response, nil
}
