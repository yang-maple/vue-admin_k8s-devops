package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
)

type persistentVolumeClaim struct{}

var PersistentVolumeClaim persistentVolumeClaim

// GetPersistentVolumeClaimList 获取 PersistentVolumeClaim 的列表
func (p *persistentVolumeClaim) GetPersistentVolumeClaimList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Claim.GetPVClaimList(params.FilterName, params.Namespace, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取持久卷声明列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取持久卷声明列表成功",
		"data": data,
	})
}

// GetPersistentVolumeClaimDetail 获取 PersistentVolumeClaim 详情
func (p *persistentVolumeClaim) GetPersistentVolumeClaimDetail(c *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `form:"persistent_volume_claim_name"`
		Namespace                 string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	detail, err := service.Claim.GetPVClaimDetail(params.Namespace, params.PersistentVolumeClaimName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷声明 " + params.PersistentVolumeClaimName + " 获取数据失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"data": detail,
		"msg":  "持久卷声明 " + params.PersistentVolumeClaimName + " 获取数据成功",
	})
}

// CreatePersistentVolumeClaim 创建PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) CreatePersistentVolumeClaim(c *gin.Context) {
	createPvc := new(struct {
		Data *service.CreateClaim `json:"data"`
	})
	_ = c.ShouldBindJSON(&createPvc)
	err := service.Claim.CreatePVClaim(createPvc.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷声明 " + createPvc.Data.Name + " 创建失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷声明 " + createPvc.Data.Name + " 创建成功",
	})
}

// DelPersistentVolumeClaim 删除 PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) DelPersistentVolumeClaim(c *gin.Context) {
	params := new(struct {
		PersistentVolumeClaimName string `json:"persistent_volume_claim_name"`
		Namespace                 string `json:"namespace"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Claim.DelPVClaim(params.Namespace, params.PersistentVolumeClaimName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷声明 " + params.PersistentVolumeClaimName + " 删除失败" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷声明 " + params.PersistentVolumeClaimName + " 删除成功",
	})
}

// UpdatePersistentVolumeClaim 更新 PersistentVolumeClaim 资源
func (p *persistentVolumeClaim) UpdatePersistentVolumeClaim(c *gin.Context) {
	params := new(struct {
		Data      *corev1.PersistentVolumeClaim `json:"data"`
		Namespace string                        `json:"namespace"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Claim.UpdatePVClaim(params.Namespace, params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷声明 " + params.Data.Name + " 更新失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷声明 " + params.Data.Name + " 更新成功",
	})
}
