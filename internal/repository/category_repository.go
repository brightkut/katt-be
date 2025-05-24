package repository

import (
	"fmt"
	"katt-be/internal/domain"
	"log"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

type CategoryRepository interface {
	Save(category *domain.Category) error
	FindById(categoryId string) (*domain.Category, error)
	FindAllByWalletId(walletId string) ([]domain.Category, error)
	Delete(category *domain.Category) error
}

func NewCategoryRepository(db *gorm.DB) categoryRepository {
	return categoryRepository{db: db}
}

func (c *categoryRepository) Save(category *domain.Category) error {
	res := c.db.Create(category)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (c *categoryRepository) FindById(categoryId string) (*domain.Category, error) {
	var category domain.Category

	res := c.db.Where("category_id = ?", categoryId).First(&category)

	fmt.Print(res)

	if res.Error != nil {
		log.Fatalf("Error get category %v", res.Error)

		return nil, res.Error
	}

	return &category, nil
}

func (c *categoryRepository) FindAllByWalletId(walletId string) ([]domain.Category, error) {
	var categories []domain.Category

	res := c.db.Where("wallet_id = ?", walletId).Find(&categories)

	if res.Error != nil {
		return nil, res.Error
	}
	return categories, nil
}

func (c *categoryRepository) Delete(category *domain.Category) error {
	res := c.db.Delete(category)

	if res.Error != nil {
		return res.Error
	}

	return nil
}
