package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID          uuid.UUID      `gorm:"column:id;type:char(36);primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"column:user_id;type:char(36)" json:"user_id"`
	Address     string         `gorm:"column:address;type:varchar(255);unique;not null" json:"address"`
	Network     string         `gorm:"column:network;type:enum('ERC20','TRC20');default:ERC20"`
	USDTBalance float64        `gorm:"column:usdt_balance" json:"usdt_balance"`
	ETHBalance  float64        `gorm:"column:eth_balance" json:"eth_balance"`
	TRXBalance  float64        `gorm:"column:trx_balance" json:"trx_balance"`
	Privatekey  string         `gorm:"column:private_key;type:varchar(255);unique" json:"private_key"`
	Publickey   string         `gorm:"column:public_key;type:varchar(255);unique" json:"public_key"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-"`
	User        *User          `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (wallet *Wallet) BeforeCreate(*gorm.DB) error {
	wallet.ID = uuid.New()
	return nil
}
