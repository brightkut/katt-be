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

func GetAllTransactionByWalletId(db *gorm.DB, dto *GetAllTransactionByWalletIdDto) (util.PaginationResult[Transaction], error) {
	query := db.Where("wallet_id = ?", dto.WalletId)

	transactions, err := util.Paginate(query, Transaction{}, int(dto.Page), int(dto.PageSize))
	if err != nil {
		return util.PaginationResult[Transaction]{}, err
	}

	return transactions, nil
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
