package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DeliverUid 捕获jwt传递的uid，并返回
func DeliverUid(c *gin.Context) *int {
	id, _ := c.Get("claims_id")
	uuid, ok := id.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取用户id失败",
			"data": nil,
		})
		return nil
	}
	return &uuid
}
