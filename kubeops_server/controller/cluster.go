package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
	"strconv"
)

var Cluster cluster

type cluster struct{}

// 默认配置文件存放路径
const configfileurl = "./static/configfile/"

// Create 创建集群
func (clt *cluster) Create(c *gin.Context) {
	//绑定参数信息
	params := new(struct {
		ClusterName string `form:"cluster_name"`
		ClusterType string `form:"cluster_type"`
	})

	_ = c.ShouldBind(&params)
	//捕获文件信息
	file, _ := c.FormFile("config_file")
	//重新dir url路径
	//uuid := c.Request.Header.Get("Uuid")
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	uid := strconv.Itoa(uuid)
	dir := configfileurl + uid

	//创建集群
	err := service.Cluster.Create(dir, params.ClusterName, params.ClusterType, uuid, file, c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "集群 " + params.ClusterName + " 创建失败:" + err.Error(),
		})
		return
	}
	//创建集群
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "集群 " + params.ClusterName + " 创建成功",
	})
}

// List 获取集群列表
func (clt *cluster) List(c *gin.Context) {
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	data, err := service.Cluster.List(uuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取集群列表成功",
		"data": data,
	})
}

// Get 获取集群信息
func (clt *cluster) Get(c *gin.Context) {
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	data, err := service.Cluster.GetClusterDetail(c.Query("cluster_name"), uuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取集群信息失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取集群信息成功",
		"data": data,
	})
}

// Change 更换集群信息
func (clt *cluster) Change(c *gin.Context) {
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	param := new(struct {
		ClusterName string `json:"cluster_name"`
	})
	_ = c.ShouldBindJSON(&param)
	err := service.Cluster.Change(param.ClusterName, uuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "切换集群 " + param.ClusterName + " 失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "已切换集群 " + param.ClusterName,
	})
}

// Delete 删除集群
func (clt *cluster) Delete(c *gin.Context) {
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	param := new(struct {
		ClusterName string `json:"cluster_name"`
	})
	_ = c.ShouldBind(&param)
	err := service.Cluster.Delete(param.ClusterName, uuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "集群 " + param.ClusterName + " 删除失败" + err.Error(),
		})
		return
	}
	//删除成功返回
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "集群 " + param.ClusterName + " 删除成功",
	})
}

// Update 更新集群信息
func (clt *cluster) Update(c *gin.Context) {
	//绑定参数信息
	params := new(struct {
		ClusterId   uint   `json:"cluster_id"`
		ClusterName string `json:"cluster_name"`
		ClusterType string `json:"cluster_type"`
	})
	_ = c.ShouldBind(&params)
	id, _ := c.Get("claims_id")
	uuid, _ := id.(int)
	//更新集群信息
	err := service.Cluster.UpdateCluster(params.ClusterId, params.ClusterName, params.ClusterType, uuid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "集群 " + params.ClusterName + " 更新失败" + err.Error(),
		})
		return
	}
	//更新成功返回
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "集群 " + params.ClusterName + " 更新成功",
	})
}
