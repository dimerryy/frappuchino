package service

import (
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type InventoryService interface {
	AddInventoryItem(item models.InventoryItem) error
	DeleteInventoryItem(id string) error
	GetInventoryItem() ([]models.InventoryItem, error)
	GetInventoryItemById(id string) (models.InventoryItem, error)
	UpdateInventoryItem(item models.InventoryItem) error
}

type inventoryService struct {
	inventoryRepo dal.InventoryRepository
}

func NewInventoryService(inventoryRepo dal.InventoryRepository) *inventoryService {
	return &inventoryService{inventoryRepo: inventoryRepo}
}

func (s *inventoryService) AddInventoryItem(item models.InventoryItem) error {
	if !IsInventoryValid(item) {
		return errors.New("invalid inventory item")
	}
	if b, _ := s.inventoryRepo.Exists(item.IngredientID); b {
		return errors.New("item already exists")
	}

	item.CreatedAt = getFormattedTime()
	item.UpdatedAt = getFormattedTime()

	return s.inventoryRepo.AddItem(item)
}

func (s *inventoryService) DeleteInventoryItem(id string) error {
	exists, err := s.inventoryRepo.Exists(id)
	if !exists {
		return errors.New("inventory item not found")
	}
	if err != nil {
		return err
	}
	return s.inventoryRepo.DeleteItem(id)
}

func (s *inventoryService) GetInventoryItem() ([]models.InventoryItem, error) {
	inventories, err := s.inventoryRepo.GetAll()
	if err != nil {
		return []models.InventoryItem{}, err
	}

	return inventories, nil
}

func (s *inventoryService) GetInventoryItemById(id string) (models.InventoryItem, error) {
	inventoryItems, err := s.inventoryRepo.GetAll()
	if err != nil {
		return models.InventoryItem{}, err
	}
	for _, inventoryItem := range inventoryItems {
		if inventoryItem.IngredientID == id {
			return inventoryItem, nil
		}
	}
	return models.InventoryItem{}, errors.New("inventory item not found")
}

func (s *inventoryService) UpdateInventoryItem(item models.InventoryItem) error {
	exists, err := s.inventoryRepo.Exists(item.IngredientID)
	if err != nil {
		return nil
	}
	if !exists {
		return errors.New("inventory item not found or you cannot change item id")
	}
	item.UpdatedAt = getFormattedTime()
	return s.inventoryRepo.UpdateItem(item)
}
