package service

import (
	"errors"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type UserService struct{}

func (s *UserService) GetUser(req database.User) (database.User, error) {
	var user database.User
	if req.ID == 0 && req.OpenID == "" && req.Phone == "" {
		return user, errors.New("必须提供 ID、OpenID 或手机号")
	}

	query := global.DB.Model(&database.User{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else if req.OpenID != "" {
		query = query.Where("openid = ?", req.OpenID)
	} else if req.Phone != "" {
		query = query.Where("phone = ?", req.Phone)
	}

	err := query.First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserService) Create(user *database.User) error {
	return global.DB.Create(user).Error
}

func (s *UserService) Update(user *database.User) error {
	return global.DB.Save(user).Error
}

func (s *UserService) Delete(where *database.User) error {
	return global.DB.Where(where).Delete(&database.User{}).Error
}
