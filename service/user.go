package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type UserService struct{}

func (userService *UserService) GetUser(req database.User) (database.User, error) {
	var user database.User
	db := global.DB
	fmt.Println(req)
	switch {
	case req.OpenID != "":
		db = db.Where("openid = ?", req.OpenID)
	case req.Nickname != "":
		db = db.Where("nickname = ?", req.Nickname)
	case req.Phone != "":
		db = db.Where("phone = ?", req.Phone)
	default:
		return database.User{}, errors.New("必须提供 openID、nickname 或 phone")
	}
	err := db.First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.User{}, err
	}
	return user, nil
}
