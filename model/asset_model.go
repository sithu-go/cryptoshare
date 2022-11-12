package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Asset struct {
	ID        uuid.UUID      `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	WalletID  uuid.UUID      `gorm:"column:wallet_id;type:char(36)" json:"wallet_id"`
	Network   string         `gorm:"column:network;type:enum('ERC20','TRC20');default:ERC20"`
	Currency  string         `gorm:"column:currency;type:enum('ETH','TRX', 'USDT');default:USDT" json:"currency"`
	Balance   float64        `gorm:"column:balance" json:"balance"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
	Wallet    *Wallet        `gorm:"foreignKey:WalletID;references:ID" json:"-"`
}

func (asset *Asset) BeforeCreate(*gorm.DB) error {
	asset.ID = uuid.New()
	return nil
}
