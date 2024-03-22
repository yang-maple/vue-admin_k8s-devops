package dao

import (
	"errors"
	"fmt"
	"kubeops/db"
	"kubeops/model"
)

// 工作流
type workflow struct{}

var Workflow workflow

type WorkflowResp struct {
	Total int               `json:"total"`
	Item  []*model.Workflow `json:"item"`
}

// GetWorkflowList 获取workflow 列表
func (w *workflow) GetWorkflowList(FilterName, Namespace string, Limit, Page int) (data *WorkflowResp, err error) {
	//定义起始页
	startSet := (Page - 1) * Limit

	//定义返回的接收变量
	var workflowList []*model.Workflow

	//查询语句 limit 限制条数 offset设置起始位置 order排序
	tx := db.GORM.Where("name like ?", "%"+FilterName+"%").Limit(Limit).Offset(startSet).Order("id desc").Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("获取workflowList失败" + err.Error())
	}
	return &WorkflowResp{
		Total: len(workflowList),
		Item:  workflowList,
	}, nil
}

// GetByIdDetail 获取详情
func (w *workflow) GetByIdDetail(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).Find(workflow)

	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("获取workflow 单条数据详情失败 没查询到该条记录" + err.Error())
	} else if tx.RowsAffected == 0 {
		return nil, errors.New("未查询到该条记录")
	}
	return workflow, nil
}

// GetById 获取单条数据
func (w *workflow) GetById(id int) (data *model.Workflow, err error) {
	var workflow model.Workflow
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	fmt.Println(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, errors.New("获取workflow 单条数据失败" + err.Error())
	}
	return &workflow, nil
}

// Add 新增数据
func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("创建workflow 数据失败" + err.Error())
	}
	return err
}

// DeleteById 删除数据
// 软删除 db.gorm.delete 实际上是更新delete_at 字段
// 硬删除 db.gorm.unscoped().delete 直接从表里删除
func (w *workflow) DeleteById(id int) (err error) {

	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return errors.New("删除workflow 数据失败" + err.Error())
	}
	return nil
}
