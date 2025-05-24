package service

import (
	"errors"
	"fmt"
	"katt-be/internal/domain"
	"katt-be/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type transactionService struct {
	walletRepository      repository.WalletRepository
	transactionRepository repository.TransactionRepository
}

type TransactionService interface {
	Create(dto *domain.CreateTransactionDto) error
	FindAllByWalletId(dto *domain.GetAllTransactionByWalletIdDto) (*domain.PaginationResult[domain.TransactionWithCategoryDTO], error)
	Delete(transactionId string) error
}

func NewTransactionService(walletRepo repository.WalletRepository, transactionRepo repository.TransactionRepository) transactionService {
	return transactionService{walletRepository: walletRepo, transactionRepository: transactionRepo}
}

func (t *transactionService) Create(dto *domain.CreateTransactionDto) error {
	transactionId, error := uuid.NewV7()

	wallet, error := t.walletRepository.FindById(dto.WalletId)

	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("wallet not found")
		}
		return error
	}

	transaction := &domain.Transaction{
		TransactionId:   transactionId.String(),
		CategoryId:      dto.CategoryId,
		WalletId:        dto.WalletId,
		TransactionType: dto.TransactionType,
		TransactionName: dto.TransactionName,
		Amount:          dto.Amount,
	}

	error = t.transactionRepository.Save(transaction)

	if error != nil {
		return error
	}

	// update amount
	if dto.TransactionType == "DEPOSIT" {
		wallet.TotalMoney += dto.Amount
	} else if dto.TransactionType == "WITHDRAW" {
		wallet.TotalMoney -= dto.Amount
	}

	error = t.walletRepository.Save(wallet)

	if error != nil {
		return error
	}

	return nil
}

func (t *transactionService) FindAllByWalletId(dto *domain.GetAllTransactionByWalletIdDto) (*domain.PaginationResult[domain.TransactionWithCategoryDTO], error) {
	transactions, error := t.transactionRepository.FindAllByWalletId(dto)

	if error != nil {
		return nil, error
	}

	return &transactions, nil
}

func (t *transactionService) Delete(transactionId string) error {

	transaction, error := t.transactionRepository.FindById(transactionId)

	if error != nil {
		return error
	}

	error = t.transactionRepository.Delete(transaction)

	if error != nil {
		return error
	}

	return nil
}
