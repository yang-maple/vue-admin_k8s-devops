package controller

import (
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"net/http"
	"strconv"
)

type workflow struct {
}

var Workflow workflow

// GetWorkflowList 获取 workflow 的list
func (w *workflow) GetWorkflowList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	_ = c.Bind(&params)
	data, err := service.Workflow.GetWorkflowList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据成功",
		"data": data,
	})

}

// GetWorkflowDetail 获取 workflow 的detail
func (w *workflow) GetWorkflowDetail(c *gin.Context) {
	params := new(struct {
		Id int `form:"id"`
	})
	_ = c.Bind(&params)
	data, err := service.Workflow.GetWorkflowDetail(params.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "获取数据详情成功",
		"data": data,
	})
}

// DeleteWorkflow 删除 workflow
func (w *workflow) DeleteWorkflow(c *gin.Context) {
	params := new(struct {
		Id int `form:"id"`
	})
	_ = c.Bind(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Workflow.DeleteById(params.Id, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

// CreateWorkflow 创建 workflow
func (w *workflow) CreateWorkflow(c *gin.Context) {
	params := new(service.Workflowcreate)
	_ = c.ShouldBindJSON(&params)
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	err := service.Workflow.CreateWorkflow(params, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "创建成功",
	})
}
