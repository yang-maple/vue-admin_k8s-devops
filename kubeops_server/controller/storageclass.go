package controller

import (
	"github.com/gin-gonic/gin"
	storagev1 "k8s.io/api/storage/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

type storageClass struct {
}

var StorageClass storageClass

// GetStorageClassList 获取存储类别列表
func (sc *storageClass) GetStorageClassList(c *gin.Context) {
	//获取实例名称,limit,page
	params := new(struct {
		FilterName string `form:"filter_name"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.StorageClass.GetStorageClassList(params.FilterName, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "获取存储类型列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "获取存储类型列表成功",
		"data":   data,
	})

}

// GetStorageClassDetail 获取存储类别详情
func (sc *storageClass) GetStorageClassDetail(c *gin.Context) {
	//获取实例名称
	param := c.Query("Name")
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.StorageClass.GetStorageClassDetail(param, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 200,
			"msg":  "存储类型 " + param + " 获取数据失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "存储类型 " + param + " 获取数据成功",
		"data":   data,
	})
}

// DelStorageClass 删除 存储类别
func (sc *storageClass) DelStorageClass(c *gin.Context) {
	//获取实例名称
	param := c.Query("Name")
	//从header获取uuid，并转为int
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.StorageClass.DelStorageClass(param, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "存储类型 " + param + " 删除失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "存储类型 " + param + " 删除成功",
		"data":   nil,
	})
}

// UpdateStorageClass 更新存储类别
func (sc *storageClass) UpdateStorageClass(c *gin.Context) {
	param := new(struct {
		Data *storagev1.StorageClass `json:"data"`
	})
	_ = c.ShouldBind(&param)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.StorageClass.UpdateStorageClass(param.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "存储类型 " + param.Data.Name + " 更新失败" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "存储类型 " + param.Data.Name + " 更新成功",
		"data": nil,
	})
}
