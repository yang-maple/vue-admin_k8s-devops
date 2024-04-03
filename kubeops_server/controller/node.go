package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
)

type node struct{}

var Node node

func (n *node) GetNodeList(c *gin.Context) {
	data, err := service.Node.GetNodeList(*DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "获取节点列表成功",
		"data": data,
	})
}

func (n *node) GetNodeDetail(c *gin.Context) {
	params := new(struct {
		NodeName string `form:"node_name"`
	})
	_ = c.ShouldBind(&params)
	data, err := service.Node.GetNodeDetail(params.NodeName, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "获取节点 " + params.NodeName + " 数据失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取节点 " + params.NodeName + " 数据成功",
		"data": data,
	})
}

// SetNodeSchedule 设置节点状态
func (n *node) SetNodeSchedule(c *gin.Context) {
	//获取节点名称以及状态
	params := new(struct {
		NodeName string `json:"node_name"`
		Status   bool   `json:"status"`
	})
	_ = c.ShouldBindJSON(&params)
	err := service.Node.SetNodeSchedule(params.NodeName, params.Status, *DeliverUid(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 4000,
			"msg":  "设置节点 " + params.NodeName + " 状态失败:" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "设置节点 " + params.NodeName + " 状态成功",
		"data": nil,
	})
}
