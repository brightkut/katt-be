package category

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateCategoryDto struct {
	WalletId     string
	CategoryName string
}

type GetAllCategoryByWalletIdDto struct {
	WalletId string
}

type Category struct {
	CategoryId   string `gorm:"primaryKey"`
	WalletId     string
	CategoryName string
	CreatedAt    time.Time // Manually define timestamp fields
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt // Enables soft delete
}

func CreateCategory(db *gorm.DB, dto *CreateCategoryDto) error {
	categoryId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	category := &Category{
		CategoryId:   categoryId.String(),
		WalletId:     dto.WalletId,
		CategoryName: dto.CategoryName,
	}

	result := db.Create(category)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllCategoryByWalletId(db *gorm.DB, dto *GetAllCategoryByWalletIdDto) []Category {
	var categories []Category
	result := db.Where("wallet_id = ?", dto.WalletId).Find(&categories)
	if result.Error != nil {
		// handle error appropriately, maybe log or panic
		return nil
	}
	return categories
}

func DeleteCategory(db *gorm.DB, categoryId string) error {
	var category Category

	result := db.First(&Category{}, categoryId)

	if result.Error != nil {
		log.Fatalf("Error get category %v", result.Error)

		return result.Error
	}

	db.Delete(&category)

	return nil
}
