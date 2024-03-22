package model

import "time"

type ClusterInfo struct {
	Id          uint `gorm:"primaryKey"`
	ClusterName string
	FileName    string
	Type        string
	Dir         string
	UserID      uint
	Status      bool
	CreatedAt   time.Time
}

func (ClusterInfo) TableName() string {
	return "clusterInfo"
}
