package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type configmap struct{}

var Configmaps configmap

// GetConfigmapList  获取 Configmap 列表
func (cm *configmap) GetConfigmapList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	_ = c.ShouldBind(&params)
	data, err := service.Configmaps.GetCmList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取配置字典列表数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取配置字典列表数据成功",
		"data": data,
	})
}

// GetConfigmapDetail   获取 Configmap 详情
func (cm *configmap) GetConfigmapDetail(c *gin.Context) {
	params := new(struct {
		ConfigmapName string `form:"configmap_name"`
		Namespace     string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	detail, err := service.Configmaps.GetCmDetail(params.Namespace, params.ConfigmapName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "配置字典 " + params.ConfigmapName + " 获取数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "配置字典 " + params.ConfigmapName + " 获取数据成功",
		"data": detail,
	})
}

// DelConfigmap    删除 Configmap 资源
func (cm *configmap) DelConfigmap(c *gin.Context) {
	params := new(struct {
		ConfigmapName string `form:"configmap_name"`
		Namespace     string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Configmaps.DelCm(params.Namespace, params.ConfigmapName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "配置字典 " + params.ConfigmapName + " 删除失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "配置字典 " + params.ConfigmapName + " 删除成功",
	})
}

// CreateConfigmap  创建 Configmap 资源
func (cm *configmap) CreateConfigmap(c *gin.Context) {
	params := new(struct {
		Data *service.CreateConfig `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Configmaps.CreateCm(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "配置字典 " + params.Data.Name + " 创建失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "配置字典 " + params.Data.Name + " 创建成功",
	})
}

// UpdateConfigmap 更新Configmap 资源
func (cm *configmap) UpdateConfigmap(c *gin.Context) {
	params := new(struct {
		Namespace string            `json:"namespace"`
		Data      *corev1.ConfigMap `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Configmaps.UpdateCm(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "配置字典 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "配置字典 " + params.Data.Name + " 更新成功",
	})
}
