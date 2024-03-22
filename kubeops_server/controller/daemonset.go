package controller

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"kubeops/service"
	"net/http"
	"strconv"
)

var DaemonSet daemonSet

type daemonSet struct {
}

// GetDaemonList 获取list
func (d *daemonSet) GetDaemonList(c *gin.Context) {

	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	data, err := service.DaemonSet.GetDsList(params.FilterName, params.Namespace, params.Limit, params.Page, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取守护进程列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取守护进程列表成功",
		"data": data,
	})
}

// GetDaemonDetail 获取demonetise 详情
func (d *daemonSet) GetDaemonDetail(c *gin.Context) {
	params := new(struct {
		DaemonName string `form:"daemon_name"`
		Namespace  string `form:"namespace"`
	})

	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	Ds, err := service.DaemonSet.GetDsDetail(params.Namespace, params.DaemonName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "守护进程 " + params.DaemonName + " 获取数据失败",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "守护进程 " + params.DaemonName + " 获取数据成功",
			"data": Ds,
		})
	}
}

// DelDaemon 删除 实例
func (d *daemonSet) DelDaemon(c *gin.Context) {
	params := new(struct {
		DaemonName string `form:"daemon_name"`
		Namespace  string `form:"namespace"`
	})
	_ = c.ShouldBind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.DaemonSet.DelDs(params.Namespace, params.DaemonName, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "守护进程 " + params.DaemonName + " 删除失败:" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "守护进程 " + params.DaemonName + " 删除成功",
		})
	}

}

// UpdateDaemon 更新实例
func (d *daemonSet) UpdateDaemon(c *gin.Context) {
	params := new(struct {
		Namespace string            `json:"namespace"`
		Data      *appsv1.DaemonSet `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.DaemonSet.UpdateDs(params.Namespace, params.Data, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "守护进程 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "守护进程 " + params.Data.Name + " 更新成功",
	})
}
