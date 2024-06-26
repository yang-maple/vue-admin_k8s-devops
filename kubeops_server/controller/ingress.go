package controller

import (
	"github.com/gin-gonic/gin"
	networkv1 "k8s.io/api/networking/v1"
	"kubeops/service"
	"net/http"
)

type ingress struct{}

var Ingress ingress

// GetIngressList 获取 Ingress 列表
func (i *ingress) GetIngressList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Ingress.GetIngList(params.FilterName, params.Namespace, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取应用路由列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取应用路由列表成功",
		"data": data,
	})
}

// GetIngressDetail  获取 Ingress 详情
func (i *ingress) GetIngressDetail(c *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
		Namespace   string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	detail, err := service.Ingress.GetIngDetail(params.Namespace, params.IngressName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由 " + params.IngressName + " 获取数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由 " + params.IngressName + " 获取数据成功",
		"data": detail,
	})
}

// DelIngress   删除 Ingress 资源
func (i *ingress) DelIngress(c *gin.Context) {
	params := new(struct {
		IngressName string `json:"ingress_name"`
		Namespace   string `json:"namespace"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Ingress.DelIng(params.Namespace, params.IngressName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由 " + params.IngressName + " 删除失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由 " + params.IngressName + " 删除成功",
	})
}

// CreateIngress  创建 Ingress 资源
func (i *ingress) CreateIngress(c *gin.Context) {
	createIngress := new(struct {
		Data *service.CreateIngress `json:"data"`
	})
	_ = c.ShouldBindJSON(&createIngress)
	err := service.Ingress.CreateIng(createIngress.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由 " + createIngress.Data.Name + " 创建失败:" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由 " + createIngress.Data.Name + " 创建成功",
	})
}

// UpdateIngress  更新 Ingress 资源
func (i *ingress) UpdateIngress(c *gin.Context) {

	params := new(struct {
		Namespace string             `json:"namespace"`
		Data      *networkv1.Ingress `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Ingress.UpdateIng(params.Namespace, params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "应用路由 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "应用路由 " + params.Data.Name + " 更新成功",
	})

}
