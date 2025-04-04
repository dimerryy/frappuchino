package service

import (
	"database/sql"
	"errors"
	"hot-coffee/internal/dal"
	"hot-coffee/models"
)

type MenuServiceInterface interface {
	AddMenuItem(item models.MenuItem) error
	GetAllMenuItems() ([]models.MenuItem, error)
	GetMenuItemById(id string) (models.MenuItem, error)
	UpdateMenu(menu models.MenuItem) error
	DeleteMenuItemById(id string) error
}

type menuService struct {
	menuRepo dal.MenuRepository
}

func NewMenuService(menuRepo dal.MenuRepository) *menuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) AddMenuItem(item models.MenuItem) error {
	if !IsMenuValid(item) {
		return errors.New("invalid menu")
	}
	exists, err := s.menuRepo.Exists(item.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("menu item already exists")
	}
	return s.menuRepo.SaveMenuItem(item)
}

func (s *menuService) GetAllMenuItems() ([]models.MenuItem, error) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return menuItems, nil
}

func (s *menuService) GetMenuItemById(id string) (models.MenuItem, error) {
	menuItems, err := s.menuRepo.GetAll()
	if err != nil {
		return models.MenuItem{}, err
	}
	for _, menuItem := range menuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}
	return models.MenuItem{}, errors.New("menu item not found")
}

func (s *menuService) UpdateMenu(menu models.MenuItem) error {
	exists, err := s.menuRepo.Exists(menu.ID)
	if err != nil {
		return err
	}
	if !exists {
		return sql.ErrNoRows
	}
	return s.menuRepo.Update(menu)
}

func (s *menuService) DeleteMenuItemById(id string) error {
	exists, err := s.menuRepo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("menu item not found")
	}
	return s.menuRepo.DeleteMenuItem(id)
}
