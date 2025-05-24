package domain

import (
	"time"

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
