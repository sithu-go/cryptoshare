package model

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	ID            uint64         `gorm:"column:id;primaryKey" json:"id"`
	Name          string         `gorm:"column:name;type:varchar(100)" json:"name"`
	WalletAddress string         `gorm:"column:wallet_address;type:varchar(255);unique" json:"wallet_address"`
	PrivateKey    string         `gorm:"column:private_key;unique" json:"private_key"`
	AddressType   string         `gorm:"column:address_type;type:enum('ERC20','TRC20');default:TRC20" json:"address_type"`
	Type          string         `gorm:"column:type" json:"type"`
	ScanRecord    string         `gorm:"column:scan_record;unique" json:"scan_record"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-"`
}

func (Bank) TableName() string {
	return "banks"
}
