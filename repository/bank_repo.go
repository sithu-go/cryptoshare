package repository

import (
	"context"
	"cryptoshare/ds"
	"cryptoshare/dto"
	"cryptoshare/model"
	"cryptoshare/utils"
	"fmt"

	"gorm.io/gorm"
)

type bankRepository struct {
	DB *gorm.DB
}

func newBankRepository(ds *ds.DataSource) *bankRepository {
	return &bankRepository{
		DB: ds.DB,
	}
}

func (r *bankRepository) Create(ctx context.Context, bank *model.Bank) error {
	return r.DB.WithContext(ctx).Debug().Create(&bank).Error
}

// in that function, we don't update private key and wallet address
func (r *bankRepository) Update(ctx context.Context, bank *model.Bank) error {
	if bank.WalletAddress != nil && bank.AddressType != nil {
		scanRecord := r.GetAddressScanRecord(*bank.AddressType, *bank.WalletAddress)
		bank.ScanRecord = &scanRecord
	}
	if bank.PrivateKey != nil {
		encrytedPrivateKey, err := utils.EncryptAES(*bank.PrivateKey)
		if err != nil {
			return err
		}
		*bank.PrivateKey = encrytedPrivateKey
	}
	return r.DB.WithContext(ctx).Debug().Updates(bank).Error
}

func (r *bankRepository) Save(bank *model.Bank) (*model.Bank, error) {
	db := r.DB.Model(&model.Bank{})
	err := db.Save(&bank).Error
	return nil, err
}

func (r *bankRepository) DeleteMany(ids string) error {
	db := r.DB.Model(&model.Bank{})
	db = db.Where(fmt.Sprintf("id in (%s)", ids))
	err := db.Delete(&model.Bank{}).Error
	return err
}

func (r *bankRepository) FindByID(id uint64) (*model.Bank, error) {
	bank := model.Bank{}
	db := r.DB.Model(&model.Bank{})
	db.Where("id", id)
	err := db.First(&bank).Error
	return &bank, err
}

func (r *bankRepository) FindAll(req *dto.RequestPayload) ([]*model.Bank, error) {
	db := r.DB.Model(&model.Bank{})
	banks := []*model.Bank{}

	db.Scopes(utils.Paginate(req.Page, req.PageSize))
	err := db.Find(&banks).Error
	return banks, err
}

func (r *bankRepository) GetAddressScanRecord(addressType string, address string) string {
	if addressType == "TRC20" {
		return fmt.Sprintf("https://tronscan.org/#/address/%s", address)
	}
	return fmt.Sprintf("https://etherscan.io/address/%s", address)
}
