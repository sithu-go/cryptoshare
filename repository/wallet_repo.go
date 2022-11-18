package repository

import (
	"cryptoshare/ds"
	"cryptoshare/model"
	"cryptoshare/service"

	"gorm.io/gorm"
)

type walletRepository struct {
	DB  *gorm.DB
	svc *service.Service
}

func newWalletRepository(ds *ds.DataSource, svc *service.Service) *walletRepository {
	return &walletRepository{
		DB:  ds.DB,
		svc: svc,
	}
}

func (r *walletRepository) Create(wallet *model.Wallet) (*model.Wallet, error) {
	db := r.DB.Model(&model.Wallet{})
	err := db.Create(&wallet).Error
	return nil, err
}
