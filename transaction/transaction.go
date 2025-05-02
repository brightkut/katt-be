package transaction

import (
	"errors"
	"fmt"
	"katt-be/util"
	"katt-be/wallet"
	"log"
	"time"

	"github.com/google/uuid"
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

func CreateTransaction(db *gorm.DB, dto *CreateTransactionDto) error {
	// Start DB transaction block
	return db.Transaction(func(tx *gorm.DB) error {
		// Generate transaction ID
		transactionId, err := uuid.NewV7()
		if err != nil {
			return err
		}

		// Find wallet first
		var w wallet.Wallet
		result := tx.First(&w, "wallet_id = ?", dto.WalletId)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return fmt.Errorf("wallet not found")
			}
			return result.Error
		}

		// Create transaction record
		transaction := &Transaction{
			TransactionId:   transactionId.String(),
			CategoryId:      dto.CategoryId,
			WalletId:        dto.WalletId,
			TransactionType: dto.TransactionType,
			TransactionName: dto.TransactionName,
			Amount:          dto.Amount,
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 🔥 Update wallet total money if DEPOSIT
		if dto.TransactionType == "DEPOSIT" {
			w.TotalMoney += dto.Amount
		} else if dto.TransactionType == "WITHDRAW" {
			w.TotalMoney -= dto.Amount
		}

		if err := tx.Save(&w).Error; err != nil {
			return err
		}

		return nil // ✅ Everything is good, commit the transaction
	})
}

func GetAllTransactionByWalletId(db *gorm.DB, dto *GetAllTransactionByWalletIdDto) (util.PaginationResult[TransactionWithCategoryDTO], error) {
	var transactions []TransactionWithCategoryDTO
	var total int64

	// Count total records
	db.Model(&Transaction{}).Where("wallet_id = ?", dto.WalletId).Count(&total)

	// Calculate pagination
	if dto.Page < 1 {
		dto.Page = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	offset := (int(dto.Page) - 1) * int(dto.PageSize)

	// Execute query with join to get category names
	err := db.Table("transactions").
		Select("transactions.*, categories.category_name").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.category_id").
		Where("transactions.wallet_id = ?", dto.WalletId).
		Order("transactions.created_at desc").
		Limit(int(dto.PageSize)).
		Offset(offset).
		Scan(&transactions).Error

	if err != nil {
		return util.PaginationResult[TransactionWithCategoryDTO]{}, err
	}

	// Create pagination result
	result := util.PaginationResult[TransactionWithCategoryDTO]{
		Page:         int(dto.Page),
		PageSize:     int(dto.PageSize),
		TotalPages:   int((total + int64(dto.PageSize) - 1) / int64(dto.PageSize)),
		TotalRecords: total,
		Data:         transactions,
	}

	return result, nil
}

func DeleteTransaction(db *gorm.DB, transactionId string) error {
	var transaction Transaction

	result := db.First(&Transaction{}, transactionId)

	if result.Error != nil {
		log.Fatalf("Error get transaction %v", result.Error)

		return result.Error
	}

	db.Delete(&transaction)

	return nil
}
