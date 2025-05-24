package repository

import (
	"gorm.io/gorm"
	"katt-be/internal/domain"
)

type walletRepository struct {
	db *gorm.DB
}

type WalletRepository interface {
	Save(wallet *domain.Wallet) error
	FindById(walletId string) (*domain.Wallet, error)
	FindByEmail(email string) (*domain.Wallet, error)
}

func NewWalletRepository(db *gorm.DB) walletRepository {
	return walletRepository{db: db}
}

func (w walletRepository) Save(wallet *domain.Wallet) error {

	res := w.db.Save(wallet)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (w *walletRepository) FindById(walletId string) (*domain.Wallet, error) {

	var wallet domain.Wallet

	res := w.db.First(&wallet, "wallet_id = ?", walletId)

	if res.Error != nil {
		return nil, res.Error
	}

	return &wallet, nil
}

func (w *walletRepository) FindByEmail(email string) (*domain.Wallet, error) {
	var wallet domain.Wallet

	res := w.db.Where("email", email).First(&wallet)

	if res.Error != nil {
		return nil, res.Error
	}

	return &wallet, nil
}
