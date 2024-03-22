package controller

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"kubeops/service"
	"kubeops/utils"
	"net/http"
	"strconv"
)

var StatefulSet statefulSet

type statefulSet struct {
}

// GetStatefulSetList 获取list
func (s *statefulSet) GetStatefulSetList(c *gin.Context) {

	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.StatefulSet.GetStatefulList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取有状态服务列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取有状态服务列表成功",
		"data": data,
	})
}

// GetStatefulDetail 获取Stateful 详情
func (s *statefulSet) GetStatefulDetail(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"stateful_set_name"`
		Namespace       string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	Ds, err := service.StatefulSet.GetStatefulDetail(params.Namespace, params.StatefulSetName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "有状态服务 " + params.StatefulSetName + " 获取数据失败",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "有状态服务 " + params.StatefulSetName + " 获取数据成功",
			"data": Ds,
		})
	}
}

// DelStatefulSet 删除 实例
func (s *statefulSet) DelStatefulSet(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `form:"stateful_set_name"`
		Namespace       string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.StatefulSet.DelStateful(params.Namespace, params.StatefulSetName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "有状态服务 " + params.StatefulSetName + " 删除失败" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "有状态服务 " + params.StatefulSetName + " 删除成功",
		})
	}

}

// UpDataStatefulSet 更新实例
func (s *statefulSet) UpDataStatefulSet(c *gin.Context) {
	params := new(struct {
		Namespace string              `json:"namespace"`
		Data      *appsv1.StatefulSet `json:"data"`
	})

	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.StatefulSet.UpdateDelStateful(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "有状态服务 " + params.Data.Name + " 更新失败：" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "有状态服务 " + params.Data.Name + " 更新成功",
	})

}

// ModifyStatefulReplicas  修改副本数
func (s *statefulSet) ModifyStatefulReplicas(c *gin.Context) {
	params := new(struct {
		StatefulSetName string `json:"stateful_set_name"`
		Namespace       string `json:"namespace"`
		Replicas        int    `json:"replicas"`
	})
	_ = c.ShouldBindJSON(&params)
	replicas := utils.Int32Ptr(int32(params.Replicas))
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.StatefulSet.ModifyStatefulReplicas(params.Namespace, params.StatefulSetName, replicas, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "有状态服务 " + params.StatefulSetName + " 更新副本数失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "有状态服务 " + params.StatefulSetName + "更新副本数成功",
	})
}
