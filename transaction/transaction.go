package transaction

import (
	"katt-be/util"
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
}

type GetAllTransactionByWalletIdDto struct {
	WalletId string
	Page     int32
	PageSize int32
}

type Transaction struct {
	TrnsactionId    string `gorm:"primaryKey"`
	WalletId        string
	CategoryId      string
	TransactionType string
	TransactionName string
	CreatedAt       time.Time // Manually define timestamp fields
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt // Enables soft delete
}

func CreateTransaction(db *gorm.DB, dto *CreateTransactionDto) error {
	transactionId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	transaction := &Transaction{
		TrnsactionId:    transactionId.String(),
		CategoryId:      dto.WalletId,
		WalletId:        dto.WalletId,
		TransactionType: dto.TransactionType,
		TransactionName: dto.TransactionName,
	}

	result := db.Create(transaction)

	if result.Error != nil {
		return result.Error
	}

	return nil
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
