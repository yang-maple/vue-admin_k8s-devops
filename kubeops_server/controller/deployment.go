package controller

import (
	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	"kubeops/service"
	"kubeops/utils"
	"net/http"
)

var Deployment deployment

type deployment struct {
}

// GetDeploylist 获取deployment列表
func (d *deployment) GetDeploylist(c *gin.Context) {

	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Deployment.GetDeploymentList(params.FilterName, params.Namespace, params.Limit, params.Page, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取无状态服务列表失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取无状态服务列表数据成功",
		"data": data,
	})
}

// ModifyDeployReplicas 修改副本数
func (d *deployment) ModifyDeployReplicas(c *gin.Context) {
	params := new(struct {
		DeployName string `json:"deploy_name"`
		Namespace  string `json:"namespace"`
		Replicas   int    `json:"replicas"`
	})

	_ = c.ShouldBindJSON(&params)
	replicas := utils.Int32Ptr(int32(params.Replicas))

	err := service.Deployment.ModifyDeployReplicas(params.Namespace, params.DeployName, replicas, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.DeployName + " 更新副本数失败:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "无状态服务 " + params.DeployName + " 更新副本数成功",
	})

}

// RestartDeploy 重启 deployment
func (d *deployment) RestartDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `json:"deploy_name"`
		Namespace  string `json:"namespace"`
	})

	_ = c.ShouldBind(&params)
	err := service.Deployment.RestartDeploy(params.Namespace, params.DeployName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.DeployName + " 重启失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "无状态服务 " + params.DeployName + " 重启成功",
		})
	}
}

// GetDeployDetail 获取 deploy 详情
func (d *deployment) GetDeployDetail(c *gin.Context) {
	params := new(struct {
		DeployName string `form:"deploy_name"`
		Namespace  string `form:"namespace"`
	})

	_ = c.ShouldBind(&params)
	deploy, err := service.Deployment.GetDeployDetail(params.Namespace, params.DeployName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.DeployName + " 获取数据失败",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "无状态服务 " + params.DeployName + " 获取数据成功",
			"data": deploy,
		})
	}
}

// CreateDeploy 创建 deploy 实例
func (d *deployment) CreateDeploy(c *gin.Context) {
	params := new(struct {
		Data *service.DeploymentCreate `json:"data"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Deployment.CreateDeploy(params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.Data.Name + " 创建失败:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "无状态服务 " + params.Data.Name + " 创建成功",
	})

}

// DelDeploy 删除 deploy 实例
func (d *deployment) DelDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `json:"deploy_name"`
		Namespace  string `json:"namespace"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Deployment.DelDeploy(params.Namespace, params.DeployName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.DeployName + " 删除失败:" + err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "无状态服务 " + params.DeployName + " 删除成功",
		})
	}
}

// UpdateDeploy 更新 deploy 实例
func (d *deployment) UpdateDeploy(c *gin.Context) {

	params := new(struct {
		Namespace string             `json:"namespace"`
		Data      *appsv1.Deployment `json:"data"`
	})

	_ = c.ShouldBindJSON(&params)
	err := service.Deployment.UpdateDeploy(params.Namespace, params.Data, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.Data.Name + " 更新失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "无状态服务 " + params.Data.Name + " 更新成功",
	})

}

// GetDeployPer 获取 每个namespace 下的deployment 实例
func (d *deployment) GetDeployPer(c *gin.Context) {
	dps, err := service.Deployment.GetDeployPer(*DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "获取无状态服务列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取无状态服务列表成功",
		"data": dps,
	})
}

// RolloutDeploy 回滚 deploy 实例
func (d *deployment) RolloutDeploy(c *gin.Context) {
	params := new(struct {
		DeployName string `form:"deploy_name"`
		Namespace  string `form:"namespace"`
	})

	_ = c.ShouldBind(&params)
	err := service.Deployment.RolloutDeploy(params.Namespace, params.DeployName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "无状态服务 " + params.DeployName + " 回滚失败" + err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 2000,
			"msg":  "无状态服务 " + params.DeployName + " 回滚成功",
		})
	}
}
