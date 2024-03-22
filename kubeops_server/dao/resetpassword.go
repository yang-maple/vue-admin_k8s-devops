package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var ResetPassword resetPassword

type resetPassword struct {
}

// VerifyIdentity 验证身份,是否有该用户的注册邮箱
func (rs *resetPassword) VerifyIdentity(email string) (err error) {
	var userInfo model.User
	tx := db.GORM.Where("email = ?", email).First(&userInfo)
	if tx.RowsAffected == 0 && tx.Error.Error() == "record not found" {
		return errors.New("邮箱未被注册")
	}
	return nil

}

// FindPasswd  用户找回并重置密码
func (rs *resetPassword) FindPasswd(email, password string) (err error) {
	var userInfo model.User
	tx := db.GORM.Model(&userInfo).Where("email = ?", email).Update("password", utils.PasswordToMd5(password))
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// ResetPasswd 重置密码
func (rs *resetPassword) ResetPasswd(Id uint, OldPasswd, NewPasswd string) (err error) {
	var userInfo model.User
	tx := db.GORM.Model(&userInfo).Where("id = ? And password = ? ", Id, utils.PasswordToMd5(OldPasswd)).First(&userInfo)
	if tx.Error != nil {
		return errors.New("旧密码输入错误" + tx.Error.Error())
	}
	tx = db.GORM.Model(&userInfo).Where("id = ?", Id).Update("password", utils.PasswordToMd5(NewPasswd))
	if tx.Error != nil {
		return errors.New("重置密码失败" + tx.Error.Error())
	}
	return nil
}
