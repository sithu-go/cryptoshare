package repository

import (
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

func (r *bankRepository) Create(bank *model.Bank) (*model.Bank, error) {
	db := r.DB.Model(&model.Bank{})
	err := db.Create(&bank).Error
	return bank, err
}

// in that function, we don't update private key and wallet address
func (r *bankRepository) Update(req *dto.UpdateBankReq) (*model.Bank, error) {
	bank, err := r.FindByID(req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		bank.Name = req.Name
	}
	if bank.WalletAddress != req.WalletAddress && req.WalletAddress != "" {
		bank.WalletAddress = req.WalletAddress
		bank.ScanRecord = r.GetAddressScanRecord(bank.AddressType, bank.WalletAddress)
	}
	if req.PrivateKey != "" {
		encrytedPrivateKey, err := utils.EncryptAES(req.PrivateKey)
		if err != nil {
			return nil, err
		}
		bank.PrivateKey = encrytedPrivateKey
	}
	if req.AddressType != "" {
		bank.AddressType = req.AddressType
		bank.ScanRecord = r.GetAddressScanRecord(bank.AddressType, bank.WalletAddress)
	}
	_, err = r.Save(bank)
	return nil, err
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

func (r *bankRepository) TransformCreateBankModel(req *dto.CreateBankReq) *model.Bank {
	scanRecord := r.GetAddressScanRecord(req.AddressType, req.WalletAddress)
	return &model.Bank{
		Name:          req.Name,
		WalletAddress: req.WalletAddress,
		PrivateKey:    req.PrivateKey,
		AddressType:   req.AddressType,
		ScanRecord:    scanRecord,
	}
}

func (r *bankRepository) GetAddressScanRecord(addressType string, address string) string {
	if addressType == "TRC20" {
		return fmt.Sprintf("https://tronscan.org/#/address/%s", address)
	}
	return fmt.Sprintf("https://etherscan.io/address/%s", address)
}
