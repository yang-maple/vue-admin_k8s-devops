package controller

import (
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"kubeops/service"
	"net/http"
)

type persistentvolume struct{}

var Persistentvolume persistentvolume

// GetPersistentVolumeList 获取PersistentVolume清单
func (p *persistentvolume) GetPersistentVolumeList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Persistenvolume.GetPvList(params.FilterName, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取持久卷列表失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取持久卷列表成功",
		"data": data,
	})
}

// GetPersistentVolumeDetail 获取 PersistentVolume 详情
func (p *persistentvolume) GetPersistentVolumeDetail(c *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `form:"persistent_volume_name"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Persistenvolume.GetPvDetail(params.PersistentVolumeName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷 " + params.PersistentVolumeName + " 获取数据失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"data": data,
	})
}

// DelPersistentVolume 删除 PersistentVolume 资源
func (p *persistentvolume) DelPersistentVolume(c *gin.Context) {
	params := new(struct {
		PersistentVolumeName string `json:"persistent_volume_name"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Persistenvolume.DelPv(params.PersistentVolumeName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷 " + params.PersistentVolumeName + " 删除失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷 " + params.PersistentVolumeName + " 删除成功",
	})
}

// CreatePersistentVolume 创建 PersistentVolume 资源
func (p *persistentvolume) CreatePersistentVolume(c *gin.Context) {
	createPv := new(struct {
		Data *service.CreatePVConfig `json:"data"`
	})
	_ = c.ShouldBindJSON(&createPv)
	err := service.Persistenvolume.CreatePv(createPv.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷 " + createPv.Data.Name + " 创建失败：" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷 " + createPv.Data.Name + " 创建成功",
	})
}

// UpdatePersistentVolume 更新PersistentVolume 资源
func (p *persistentvolume) UpdatePersistentVolume(c *gin.Context) {
	params := new(struct {
		Data *corev1.PersistentVolume `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Persistenvolume.UpdatePv(params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "持久卷 " + params.Data.Name + " 更新失败：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "持久卷 " + params.Data.Name + " 更新成功",
	})

}
