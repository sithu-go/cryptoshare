package ds

import (
	"cryptoshare/model"
	"cryptoshare/utils"
	"log"
)

func AddDefaultAdmin() {
	tb := DB.Model(&model.Admin{})
	var count int64
	tb.Count(&count)
	if count != 0 {
		return
	}
	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		log.Panic(err)
	}
	admin := model.Admin{
		Name:     utils.NewString("Admin"),
		Username: utils.NewString("admin"),
		Email:    utils.NewString("admin@gmail.com"),
		Password: &hashedPassword,
		IP:       utils.NewString("12.12.12.12"),
	}

	if err := tb.Create(&admin).Error; err != nil {
		log.Panic(err)
		return
	}
}
