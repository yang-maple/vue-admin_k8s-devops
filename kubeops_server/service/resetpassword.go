package service

import (
	"errors"
	"kubeops/dao"
	"kubeops/utils"
)

var ResetPassword resetPassword

type resetPassword struct{}

// VerifyIdentity 验证用户身份
func (rs *resetPassword) VerifyIdentity(email string) (err error) {
	// 判断用户是否存在
	if err := dao.ResetPassword.VerifyIdentity(email); err != nil {
		utils.Logger.Error("The email address is not registered")
		return err
	}
	//组装数据，生成随机数
	verity := emailVerity{
		Email: email + "_rsp",
		Value: utils.GenRandNum(6),
	}
	// 将随机数缓存至redis
	err = dao.RdbValue.SetValue(verity.Email, verity.Value, 90)
	if err != nil {
		utils.Logger.Error("Failed to cache the random number to Redis" + err.Error())
		return err
	}
	// 随机数验证码发送邮件
	err = utils.Emails(email, utils.FormatEmailBody("view/verfityidentify.html", verity), "KubeOps找回密码")
	if err != nil {
		utils.Logger.Error("Failed to send an email to " + email + " :" + err.Error())
		return err
	}
	utils.Logger.Info("The verification code has been sent to " + email)
	return nil
}

// FindPasswd  找回并重置密码
func (rs *resetPassword) FindPasswd(email, verifyCode, password string) (err error) {
	//读取缓存中的验证码
	value, err := dao.RdbValue.GetValue(email + "_rsp")
	if err != nil {
		utils.Logger.Error("Failed to read cached nonce in Redis,reason: " + err.Error())
		return err
	}
	if value != verifyCode {
		utils.Logger.Error("The verification code is incorrect")
		return errors.New("验证码错误")
	}
	//重置密码
	if err := dao.ResetPassword.FindPasswd(email, password); err != nil {
		utils.Logger.Error("Failed to reset the password" + err.Error())
		return err
	}
	//重置成功后删除缓存中的验证码
	_ = dao.RdbValue.DelValue(email + "_rsp")
	utils.Logger.Info("The password has been reset successfully")
	return nil

}

// ResetPasswd 重置密码
func (rs *resetPassword) ResetPasswd(Id uint, OldPasswd, NewPasswd string) (err error) {
	//数据库验证密码
	err = dao.ResetPassword.ResetPasswd(Id, OldPasswd, NewPasswd)
	if err != nil {
		utils.Logger.Error("Failed to reset the password" + err.Error())
		return err
	}
	utils.Logger.Info("The password has been reset successfully")
	return nil
}
