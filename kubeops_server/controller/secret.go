package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type secret struct{}

var Secrets secret

// GetSecretList   获取 Secret 列表
func (s *secret) GetSecretList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.Secrets.GetSecretList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取保密字典列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取保密字典列表成功",
		"data": data,
	})
}

// GetSecretDetail    获取 Secret 详情
func (s *secret) GetSecretDetail(c *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	detail, err := service.Secrets.GetSecretDetail(params.Namespace, params.SecretName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "保密字典 " + params.SecretName + " 获取数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "保密字典 " + params.SecretName + " 获取数据成功",
		"data": detail,
	})
}

// DelSecret    删除 Secret 资源
func (s *secret) DelSecret(c *gin.Context) {
	params := new(struct {
		SecretName string `form:"secret_name"`
		Namespace  string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Secrets.DelSecret(params.Namespace, params.SecretName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "保密字典 " + params.SecretName + " 删除失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "保密字典 " + params.SecretName + " 删除成功",
	})
}

// CreateSecret   创建 Secret 资源
func (s *secret) CreateSecret(c *gin.Context) {
	params := new(struct {
		Data *service.CreateSecret `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Secrets.CreateSecret(params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "保密字典 " + params.Data.Name + " 创建失败" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "保密字典 " + params.Data.Name + " 创建成功",
	})
}

// UpdateSecret   更新 Secret 资源
func (s *secret) UpdateSecret(c *gin.Context) {
	params := new(struct {
		Namespace string         `json:"namespace"`
		Data      *corev1.Secret `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Secrets.UpdateSecret(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "保密字典 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "保密字典 " + params.Data.Name + " 更新成功",
	})
}
