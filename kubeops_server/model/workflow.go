package model

import (
	"time"
)

type Workflow struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	//gorm:"column:type" 用于声明 mysql 中的表字段名
	Type      string     `json:"service_type" gorm:"column:type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Workflow) TableName() string {
	return "workflow"
}
