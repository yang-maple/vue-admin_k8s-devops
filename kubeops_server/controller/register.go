package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
)

var Register register

type register struct{}

// RegisterUser 注册用户
func (r *register) RegisterUser(c *gin.Context) {
	var registerInfo service.RegisterInfo
	_ = c.ShouldBindJSON(&registerInfo)
	err := service.Register.RegisterUser(&registerInfo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  " 注册失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "用户 " + registerInfo.Username + " 注册成功",
	})
}

// SendEmail 注册用户时发送邮件
func (r *register) SendEmail(c *gin.Context) {
	params := new(struct {
		Email string `json:"email"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Register.SendEmail(params.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "验证码已发送至邮箱",
	})
}
