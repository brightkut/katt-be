package wallet

import (
	"time"

	"github.com/google/uuid"
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

func CreateWallet(db *gorm.DB, dto *CreateWalletDto) error {
	walletId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	wallet := &Wallet{
		WalletId:   walletId.String(),
		Email:      dto.Email,
		TotalMoney: dto.TotalMoney,
	}

	result := db.Create(wallet)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetWalletByEmail(db *gorm.DB, dto *GetWalletByEmailDto) *Wallet {
	var wallet Wallet

	result := db.Where("email", dto.Email).First(&wallet)

	if result.Error != nil {
		return nil
	}

	return &wallet
}
