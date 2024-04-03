package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
)

type namespace struct{}

var Namespace namespace

func (n *namespace) GetNsList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)

	data, err := service.Namespace.GetNsList(params.FilterName, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取命名空间列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取命名空间列表成功",
		"data": data,
	})
}

func (n *namespace) GetNsDetail(c *gin.Context) {
	params := new(struct {
		NamespaceName string `form:"namespace_name"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Namespace.GetNsDetail(params.NamespaceName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取  " + params.NamespaceName + " 失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取  " + params.NamespaceName + "数据成功",
		"data": data,
	})
}

func (n *namespace) DelNs(c *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
	})
	_ = c.ShouldBind(&params)
	err := service.Namespace.DelNs(params.NamespaceName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "命名空间 " + params.NamespaceName + " 删除失败" + " " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "命名空间 " + params.NamespaceName + " 删除成功",
	})
}

func (n *namespace) CreateNs(c *gin.Context) {
	params := new(struct {
		NamespaceName string `json:"namespace_name"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Namespace.CreateNs(params.NamespaceName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "命名空间 " + params.NamespaceName + " 创建失败" + " " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "命名空间 " + params.NamespaceName + " 创建成功",
	})
}

func (n *namespace) UpdateNs(c *gin.Context) {
	params := new(struct {
		Data *corev1.Namespace `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Namespace.UpdateNs(params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "命名空间 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "命名空间 " + params.Data.Name + " 更新成功",
	})

}
