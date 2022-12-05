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

type adminRepository struct {
	DB *gorm.DB
}

func newAdminRepository(ds *ds.DataSource) *adminRepository {
	return &adminRepository{
		DB: ds.DB,
	}
}

func (r *adminRepository) FindOrByField(field1, field2, value string) (*model.Admin, error) {
	db := r.DB.Model(&model.Admin{})
	admin := model.Admin{}
	err := db.First(&admin, fmt.Sprintf("%s = ? OR %s = ?", field1, field2), value, value).Error
	return &admin, err
}

func (r *adminRepository) FindByField(field, value string) (*model.Admin, error) {
	db := r.DB.Model(&model.Admin{})
	admin := model.Admin{}
	err := db.First(&admin, fmt.Sprintf("%s = ?", field), value).Error
	return &admin, err
}

func (r *adminRepository) UpdateByFields(ctx context.Context, updateFields *model.UpdateFields) (*model.Admin, error) {
	db := r.DB.WithContext(ctx).Debug().Model(&model.Admin{})
	db.Where(updateFields.Field, updateFields.Value)
	err := db.Updates(updateFields.Data).Error
	return nil, err
}

func (r *adminRepository) Create(ctx context.Context, admin *model.Admin) error {
	return r.DB.WithContext(ctx).Debug().Model(&model.Admin{}).Create(&admin).Error
}

func (r *adminRepository) List(ctx context.Context, req *dto.PageReq) ([]*model.Admin, int64, error) {
	tb := r.DB.WithContext(ctx).Debug().Model(&model.Admin{})
	var total int64
	tb.Count(&total)
	tb.Scopes(utils.Paginate(req.Page, req.PageSize))
	admins := make([]*model.Admin, 0)
	return admins, total, tb.Find(&admins).Error
}

func (r *adminRepository) Update(ctx context.Context, admin *model.Admin) error {
	return r.DB.WithContext(ctx).Debug().Where("id", admin.ID).UpdateColumns(admin).Error
}

func (r *adminRepository) DeleteMany(ctx context.Context, ids string) error {
	return r.DB.WithContext(ctx).Debug().Delete(&model.Admin{}, fmt.Sprintf("id in (%s)", ids)).Error
}

func (r *adminRepository) RecoverAdmins(ctx context.Context, ids string) error {
	return r.DB.WithContext(ctx).Debug().Unscoped().Model(&model.Admin{}).Where("id IN (?)", ids).Update("deleted_at", nil).Error
}
