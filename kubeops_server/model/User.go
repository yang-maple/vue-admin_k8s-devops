package model

import (
	"time"
)

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "editor"
	UserRoleGuest UserRole = "visitor"
)

type User struct {
	Id           uint `gorm:"primaryKey"`
	Username     string
	Password     string
	Email        string
	Roles        UserRole
	Avatar       string
	Status       int
	CreatedAt    time.Time
	ClustersInfo []ClusterInfo `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "user"
}
