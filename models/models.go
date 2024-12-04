package models

import "gorm.io/gorm"

type Wallet struct {
	UUID    string  `gorm:"primaryKey;column:id"`
	Balance float64 `gorm:"not null"`
}

type WalletReq struct {
	WalletID      string  `json:"walletId" validate:"required,uuid"`
	OperationType string  `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

func (w *Wallet) UpdateBalance(tx *gorm.DB, operationType string, amount float64) error {
	if operationType == "WITHDRAW" && w.Balance < amount {
		return gorm.ErrInvalidData
	}
	if operationType == "DEPOSIT" {
		w.Balance += amount
	} else if operationType == "WITHDRAW" {
		w.Balance -= amount
	}
	return tx.Save(w).Error
}
