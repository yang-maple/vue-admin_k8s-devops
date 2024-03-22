package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type services struct{}

var Services services

// GetServiceList 获取 services 列表
func (s *services) GetServiceList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Services.GetSvcList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "获取服务列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取服务列表成功",
		"data": data,
	})
}

// GetServiceDetail 获取 services 详情
func (s *services) GetServiceDetail(c *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	detail, err := service.Services.GetSvcDetail(params.Namespace, params.ServiceName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "服务 " + params.ServiceName + " 获取数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "服务 " + params.ServiceName + " 获取数据成功",
		"data": detail,
	})
}

// DelServices  删除 services 资源
func (s *services) DelServices(c *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Services.DelSvc(params.Namespace, params.ServiceName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "服务 " + params.ServiceName + " 删除失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "服务 " + params.ServiceName + " 删除成功",
	})
}

// CreateService 创建 Services 资源
func (s *services) CreateService(c *gin.Context) {
	createSvc := new(struct {
		Data *service.CreateService `json:"data"`
	})
	_ = c.ShouldBindJSON(&createSvc)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Services.CreateSvc(createSvc.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "服务 " + createSvc.Data.Name + " 创建失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "服务 " + createSvc.Data.Name + " 创建成功",
	})
}

// UpdateService 更新 Services 资源
func (s *services) UpdateService(c *gin.Context) {

	params := new(struct {
		Namespace string          `json:"namespace"`
		Data      *corev1.Service `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Services.UpdateSvc(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "服务 " + params.Data.Name + " 更新失败：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "服务 " + params.Data.Name + " 更新成功",
	})
}
