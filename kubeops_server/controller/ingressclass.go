package controller

import (
	"github.com/gin-gonic/gin"
	networkv1 "k8s.io/api/networking/v1"
	"kubeops/service"
	"net/http"
)

var IngressClass ingressClass

type ingressClass struct {
}

// GetIngressClassList list 获取列表信息
func (ic *ingressClass) GetIngressClassList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	//从header获取uuid，并转为int
	data, err := service.IngressClass.GetIngressClass(params.FilterName, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取应用路由类型列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取应用路由类型列表成功",
		"data": data,
	})

}

// GetIngressClassDetail 获取详情信息
func (ic *ingressClass) GetIngressClassDetail(c *gin.Context) {
	param := c.Query("Name")
	//从header获取uuid，并转为int
	data, err := service.IngressClass.GetIngressClassDetail(param, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由类型 " + param + " 获取数据失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由类型 " + param + " 获取数据成功",
		"data": data,
	})

}

// DelIngressClass 删除ingressClass
func (ic *ingressClass) DelIngressClass(c *gin.Context) {
	param := c.Query("Name")
	//从header获取uuid，并转为int
	err := service.IngressClass.DelIngressClass(param, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由类型 " + param + " 删除失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由类型 " + param + " 删除成功",
	})
}

// UpdateIngressClass 更新ingressClass
func (ic *ingressClass) UpdateIngressClass(c *gin.Context) {
	param := new(struct {
		Data *networkv1.IngressClass `json:"data"`
	})
	_ = c.ShouldBindJSON(&param)
	err := service.IngressClass.UpdateIngressClass(param.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由类型 " + param.Data.Name + " 更新失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由类型 " + param.Data.Name + " 更新成功",
		"data": nil,
	})
}
