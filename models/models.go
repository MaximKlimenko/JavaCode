package models

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type Wallet struct {
	ID      string  `gorm:"primaryKey;column:id"`
	Balance float64 `gorm:"not null"`
}

type WalletReq struct {
	WalletID      string  `json:"walletId" validate:"required,id"`
	OperationType string  `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

func (w *Wallet) UpdateBalance(tx *gorm.DB, operationType string, amount float64) error {
	switch operationType {
	case "DEPOSIT":
		w.Balance += amount
	case "WITHDRAW":
		if w.Balance < amount {
			return gorm.ErrInvalidData // Недостаточно средств
		}
		w.Balance -= amount
	default:
		return errors.New("invalid operation type") // Некорректная операция
	}

	// Сохраняем изменения в базе
	if err := tx.Save(w).Error; err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	return nil
}
