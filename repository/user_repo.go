package repository

import (
	"cryptoshare/ds"
	"cryptoshare/model"
	"fmt"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func newUserRepository(ds *ds.DataSource) *userRepository {
	return &userRepository{
		DB: ds.DB,
	}
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	db := r.DB.Model(&model.User{})
	err := db.Create(&user).Error
	return user, err
}

func (r *userRepository) FindByField(field, value string) (*model.User, error) {
	db := r.DB.Model(&model.User{})
	user := model.User{}
	err := db.First(&user, fmt.Sprintf("%s = ?", field), value).Error
	return &user, err
}

func (r *userRepository) UpdateByFields(updateFields *model.UpdateFields) (*model.User, error) {
	db := r.DB.Model(&model.User{})
	db.Where(updateFields.Field, updateFields.Value)
	err := db.Updates(updateFields.Data).Error
	return nil, err
}
