package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var Register register

type register struct{}

func (r *register) RegisterUser(username, password, email string) (err error) {
	var userInfo model.User
	//判断用户是否存在
	tx := db.GORM.Where("username = ?", username).First(&userInfo)
	if tx.Error != nil && tx.Error.Error() != "record not found" || tx.RowsAffected != 0 {
		return errors.New("用户名已存在")
	}
	//组装新用户信息
	newUser := model.User{
		Username: username,
		Password: utils.PasswordToMd5(password),
		Email:    email,
		Status:   1,
		//默认用户角色为editor
		Roles: model.UserRoleUser,
		//默认用户头像
		Avatar: "https://ts2.cn.mm.bing.net/th?id=OIP-C.jHUH4s7TQ48X_B-1iozuJgHaHa&w=250&h=250&c=8&rs=1&qlt=90&o=6&dpr=1.3&pid=3.1&rm=2",
	}
	//注册用户
	tx = db.GORM.Create(&newUser)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// SendEmail 验证邮箱是否已被注册
func (r *register) SendEmail(email string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("email = ?", email).First(&userInfo)
	if tx.RowsAffected != 0 {
		return errors.New("邮箱已被注册")
	}
	return nil
}
