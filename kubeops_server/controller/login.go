package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"kubeops/service"
	"net/http"
)

type login struct{}

var Login login

var store = base64Captcha.DefaultMemStore

// CaptchaImage  生成验证码,并返回图片url
func (l *login) CaptchaImage(c *gin.Context) {
	//定义图片格式为 string
	driverConfig := base64Captcha.DriverString{
		Height:          43,                                     //高度
		Width:           130,                                    //宽度
		NoiseCount:      0,                                      //数字干扰项，数字越大，
		ShowLineOptions: 2,                                      //竖线干扰项
		Length:          6,                                      //文本长度
		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm", //文本来源
		BgColor: &color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 1,
		},
		Fonts: []string{"wqy-microhei.ttc"}, // 文本样式可选
		// "3Dumb.ttf",
		// "ApothecaryFont.ttf",
		// "Comismsh.ttf",
		// "DENNEthree-dee.ttf",
		// "DeborahFancyDress.ttf",
		// "Flim-Flam.ttf",
		// "RitaSmith.ttf",
		// "actionj.ttf",
		// "chromohv.ttf",
		// "wqy-microhei.ttc",
	}
	//生成图片
	captcha := base64Captcha.NewCaptcha(driverConfig.ConvertFonts(), store)
	//获取图片id,url
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      2000,
		"imageUrl":  b64s,
		"captchaId": id,
		"msg":       "success",
	})
}

// VerifyInfo 验证登陆信息
func (l *login) VerifyInfo(c *gin.Context) {
	params := new(service.LoginInfo)
	_ = c.ShouldBindJSON(&params)
	//判断验证码是否正确 不正确直接返回
	if !store.Verify(params.CaptchaId, params.VerifyValue, true) {
		c.JSON(http.StatusOK, gin.H{
			"code":  4000,
			"msg":   "验证码错误",
			"token": "",
		})
		return
	}
	// 验证码验证成功后，判断用户名密码
	token, err := service.Login.VerifyUserInfo(params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":  4000,
			"msg":   err.Error(),
			"token": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":  2000,
		"msg":   params.Username + "登陆成功",
		"token": token,
	})
}

func (l *login) getUserInfo(c *gin.Context) {
	params := new(struct {
		Token string `form:"token"`
	})
	_ = c.ShouldBind(&params)
	info, err := service.Login.GetUserInfo(params.Token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取用户信息成功",
		"data": info,
	})

}
