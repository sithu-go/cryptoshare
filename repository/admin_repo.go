package repository

import (
	"cryptoshare/ds"
	"cryptoshare/model"
	"fmt"

	"gorm.io/gorm"
)

type adminRepository struct {
	DB *gorm.DB
}

func newAdminRepository(ds *ds.DataSource) *adminRepository {
	return &adminRepository{
		DB: ds.DB,
	}
}

func (r *adminRepository) FindByField(field, value string) (*model.Admin, error) {
	db := r.DB.Model(&model.Bank{})
	admin := model.Admin{}
	err := db.First(&admin, fmt.Sprintf("%s = ?", field), value).Error
	return &admin, err
}
