package service

import (
	"fmt"
	"katt-be/internal/domain"
	"katt-be/internal/repository"
	"log"

	"github.com/google/uuid"
)

type walletService struct {
	walletRepository repository.WalletRepository
}

type WalletService interface {
	Create(dto *domain.CreateWalletDto) error
	GetByEmail(dto *domain.GetWalletByEmailDto) *domain.Wallet
}

func NewWalletService(walletRepository repository.WalletRepository) walletService {
	return walletService{walletRepository: walletRepository}
}

func (w *walletService) Create(dto *domain.CreateWalletDto) error {
	walletId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	wallet := &domain.Wallet{
		WalletId:   walletId.String(),
		Email:      dto.Email,
		TotalMoney: dto.TotalMoney,
	}

	return w.walletRepository.Save(wallet)
}

func (w *walletService) GetByEmail(dto *domain.GetWalletByEmailDto) *domain.Wallet {
	wallet, err := w.walletRepository.FindByEmail(dto.Email)

	if err != nil {
		log.Fatal(fmt.Printf("Error can't get wallet by email: %v", err))
		return nil
	}

	return wallet
}
