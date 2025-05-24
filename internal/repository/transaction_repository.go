package repository

import (
	"katt-be/internal/domain"
	"log"

	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

type TransactionRepository interface {
	Save(transaction *domain.Transaction) error
	Delete(transaction *domain.Transaction) error
	FindById(transactionId string) (*domain.Transaction, error)
	FindAllByWalletId(dto *domain.GetAllTransactionByWalletIdDto) (domain.PaginationResult[domain.TransactionWithCategoryDTO], error)
}

func NewTransactionRepository(db *gorm.DB) transactionRepository {
	return transactionRepository{db: db}
}

func (t *transactionRepository) Save(transaction *domain.Transaction) error {
	res := t.db.Create(transaction)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (t *transactionRepository) Delete(transaction *domain.Transaction) error {
	res := t.db.Delete(transaction)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (t *transactionRepository) FindById(transactionId string) (*domain.Transaction, error) {
	var transaction domain.Transaction

	res := t.db.First(&domain.Transaction{}, transactionId)

	if res.Error != nil {
		log.Fatalf("Error get transaction %v", res.Error)

		return nil, res.Error
	}

	return &transaction, nil
}

func (t *transactionRepository) FindAllByWalletId(dto *domain.GetAllTransactionByWalletIdDto) (domain.PaginationResult[domain.TransactionWithCategoryDTO], error) {
	var transactions []domain.TransactionWithCategoryDTO
	var total int64

	// Count total records
	t.db.Model(&domain.Transaction{}).Where("wallet_id = ?", dto.WalletId).Count(&total)

	// Calculate pagination
	if dto.Page < 1 {
		dto.Page = 1
	}
	if dto.PageSize <= 0 {
		dto.PageSize = 10
	}

	offset := (int(dto.Page) - 1) * int(dto.PageSize)

	// Execute query with join to get category names
	err := t.db.Table("transactions").
		Select("transactions.*, categories.category_name").
		Joins("LEFT JOIN categories ON transactions.category_id = categories.category_id").
		Where("transactions.wallet_id = ?", dto.WalletId).
		Order("transactions.created_at desc").
		Limit(int(dto.PageSize)).
		Offset(offset).
		Scan(&transactions).Error

	if err != nil {
		return domain.PaginationResult[domain.TransactionWithCategoryDTO]{}, err
	}

	// Create pagination result
	result := domain.PaginationResult[domain.TransactionWithCategoryDTO]{
		Page:         int(dto.Page),
		PageSize:     int(dto.PageSize),
		TotalPages:   int((total + int64(dto.PageSize) - 1) / int64(dto.PageSize)),
		TotalRecords: total,
		Data:         transactions,
	}

	return result, nil
}
