package model

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	ID         *uint64        `gorm:"column:id;primaryKey" json:"id"`
	Name       *string        `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Username   *string        `gorm:"column:username;type:varchar(100);unique;not null" json:"username"`
	Email      *string        `gorm:"column:email;type:varchar(100);unique;not null" json:"email"`
	Password   *string        `gorm:"column:password;type:varchar(255);not null" json:"-"`
	IP         *string        `gorm:"column:ip;type:varchar(20)" json:"ip"`
	OTPEnabled *bool          `gorm:"column:otp_enabled;default:false;not null" json:"otp_enabled"`
	OTPSecret  *string        `gorm:"column:otp_secret" json:"otp_secret"`
	OTPAuthURL *string        `gorm:"column:otp_auth_url;default:false;not null" json:"otp_auth_url"`
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-"`
	// OTPVerified bool           `gorm:"column:otp_verified;default:false;not null" json:"otp_verified"`
}
