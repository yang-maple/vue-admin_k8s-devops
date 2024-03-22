package dao

import (
	"errors"
	"fmt"
	"kubeops/db"
	"kubeops/model"
	"kubeops/utils"
)

var Login login

type login struct {
}

// VerifyUser 验证用户名和密码 以及  账号状态
func (l *login) VerifyUser(username, password string) (id int, err error) {
	var userInfo model.User
	//验证用户身份
	tx := db.GORM.Where("username = ?", username).Where("password = ?", utils.PasswordToMd5(password)).First(&userInfo)
	//返回查询错误
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return -1, tx.Error
	}
	//返回记录不存在
	if tx.RecordNotFound() {
		return -1, errors.New("用户不存在")
	}

	// 如果你的库实现中 tx.RecordNotFound() 并不能准确反映查询结果为空的情况，
	// 或者你需要同时处理 RowsAffected == 0 的其他逻辑，那么保留下面这行：
	// if tx.RowsAffected == 0 {
	//     return nil, errors.New("用户不存在")
	// }
	//判断用户状态
	if userInfo.Status == 0 {
		return -1, errors.New("用户已被禁用")
	}
	return int(userInfo.Id), nil
}

// GetUserInfo 获取用户信息(用户名，角色，集群url, 集群名称)
func (l *login) GetUserInfo(id int) (*model.User, error) {
	var userInfo model.User
	// 获取用户信息
	tx := db.GORM.Where("id = ?", id).First(&userInfo)
	//返回查询错误
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return nil, tx.Error
	}
	//返回记录不存在
	if tx.RecordNotFound() {
		return nil, errors.New("用户不存在")
	}

	// 如果你的库实现中 tx.RecordNotFound() 并不能准确反映查询结果为空的情况，
	// 或者你需要同时处理 RowsAffected == 0 的其他逻辑，那么保留下面这行：
	// if tx.RowsAffected == 0 {
	//     return nil, errors.New("用户不存在")
	// }
	return &userInfo, nil
}

// GetUserCluster 获取用户集群信息
func (l *login) GetUserCluster(id int) (url, name string, err error) {
	var cluster model.ClusterInfo
	//判断用户是否有集群
	tx := db.GORM.Where("user_id =? And status = ?", id, 1).First(&cluster)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		return "", "", tx.Error
	}
	//没有集群返回 nil 值
	if tx.RecordNotFound() {
		return "", "", nil
	}
	// 有集群返回集群信息
	//return int(userInfo.Id), cluster.Dir + "/" + cluster.ClusterName + "_" + cluster.Type + ".conf", cluster.ClusterName, nil
	return fmt.Sprintf("%s/%s_%s.conf", cluster.Dir, cluster.ClusterName, cluster.Type), cluster.ClusterName, err
}
