package domain

import (
	"time"

	"gorm.io/gorm"
)

type CreateWalletDto struct {
	Email      string `json:"email"`
	TotalMoney float64
}

type GetWalletByEmailDto struct {
	Email string `json:"email"`
}

type Wallet struct {
	WalletId   string `gorm:"primaryKey"`
	Email      string `json:"email" gorm:"unique"`
	TotalMoney float64
	CreatedAt  time.Time // Manually define timestamp fields
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt // Enables soft delete
}
