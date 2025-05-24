package domain

import (
	"time"

	"gorm.io/gorm"
)

type CreateTransactionDto struct {
	WalletId        string
	CategoryId      string
	TransactionType string
	TransactionName string
	Amount          float64
}

type GetAllTransactionByWalletIdDto struct {
	WalletId string
	Page     int32
	PageSize int32
}

type TransactionWithCategoryDTO struct {
	TransactionId   string    `json:"transactionId"`
	WalletId        string    `json:"walletId"`
	CategoryId      string    `json:"categoryId"`
	CategoryName    string    `json:"categoryName"` // Added category name
	TransactionType string    `json:"transactionType"`
	TransactionName string    `json:"transactionName"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Transaction struct {
	TransactionId   string `gorm:"primaryKey"`
	WalletId        string
	CategoryId      string
	TransactionType string
	TransactionName string
	Amount          float64
	CreatedAt       time.Time // Manually define timestamp fields
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt // Enables soft delete
}
