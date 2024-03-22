package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
	"strconv"
)

var HomePage homepage

type homepage struct {
}

// GetHomepage  获取首页node数据
func (h *homepage) GetHomepage(c *gin.Context) {
	userid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Homepage.GetHomepage(userid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}
