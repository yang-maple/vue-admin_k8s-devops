package dao

import (
	"errors"
	"kubeops/db"
	"kubeops/model"
)

var Cluster cluster

type cluster struct{}

// AddCluster 添加集群
func (clt *cluster) AddCluster(clusterName, filename, clusterType, dir string, uuid int) error {
	var cluster model.ClusterInfo
	//判断集群是否存在
	tx := db.GORM.Where("cluster_name = ?", clusterName).First(&cluster)
	if tx.RowsAffected != 0 {
		return errors.New("cluster name already exists")
	}
	//组装新集群信息
	newCluster := model.ClusterInfo{
		ClusterName: clusterName,
		FileName:    filename,
		Dir:         dir,
		Type:        clusterType,
		UserID:      uint(uuid),
		Status:      false,
	}
	//注册集群
	tx = db.GORM.Create(&newCluster)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// ListCluster List 获取集群列表
func (clt *cluster) ListCluster(uuid int) (clusters *[]model.ClusterInfo, err error) {
	var cluster []model.ClusterInfo
	tx := db.GORM.Where("user_id = ?", uuid).Find(&cluster)
	if tx.RowsAffected < 1 {
		return nil, errors.New("the user did not add cluster information")
	}
	return &cluster, nil
}

// ChangeCluster  更换集群信息
func (clt *cluster) ChangeCluster(clusterName string, uuid int) (err error) {
	//旧集群状态设置为false
	db.GORM.Model(&model.ClusterInfo{}).Where("user_id = ?", uuid).Where("status = ?", true).Update("status", false)
	//更新新集群状态为true
	tx := db.GORM.Model(&model.ClusterInfo{}).Where("cluster_name = ?  AND user_id = ?", clusterName, uuid).Update("status", true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// DeleteCluster 删除集群
func (clt *cluster) DeleteCluster(name string, userid int) (err error) {
	//定义模型
	var cluster model.ClusterInfo
	if db.GORM.Where("cluster_name =? AND user_id =?", name, userid).Delete(&cluster).RowsAffected == 0 {
		return errors.New("the cluster does not exist")
	}
	return nil
}

// GetClusterDir 获取集群dir
func (clt *cluster) GetClusterDir(name string, userid int) (dir string, err error) {
	var cluster model.ClusterInfo
	tx := db.GORM.Where("cluster_name =? AND user_id =?", name, userid).Find(&cluster)
	if tx.RowsAffected == 0 {
		return "", errors.New("the cluster does not exist")
	}
	return cluster.Dir + "/" + cluster.ClusterName + "_" + cluster.Type + ".conf", nil
}

// GetClusterDetail 获取集群详情
func (clt *cluster) GetClusterDetail(name string, userid int) (clusters *model.ClusterInfo, err error) {
	var cluster model.ClusterInfo
	tx := db.GORM.Where("cluster_name =? AND user_id =?", name, userid).Find(&cluster)
	if tx.RowsAffected == 0 {
		return nil, errors.New("the cluster does not exist")
	}
	return &cluster, nil
}

// UpdateCluster 更新集群信息
func (clt *cluster) UpdateCluster(Id uint, NewName, NewType string, uuid int) (dir, name string, err error) {
	var cluster model.ClusterInfo
	//先根据UID 找到用户集群信息，然后判断新集群名字是否存在 如果存在则返回错误，如果不存在则更新集群信息
	tx := db.GORM.Where("cluster_name =? AND user_id =?", NewName, uuid).First(&cluster)
	if tx.RowsAffected != 0 {
		return "", "", errors.New("the cluster already exists")
	}
	//根据ID 找到集群信息
	tx = db.GORM.Where("id = ?", Id).First(&cluster)
	//获取旧集群名称
	name = cluster.Dir + "/" + cluster.ClusterName + "_" + cluster.Type + ".conf"
	//获取旧集群路径
	dir = cluster.Dir
	// 更新集群字段
	cluster.ClusterName = NewName
	cluster.Type = NewType
	tx = db.GORM.Save(&cluster)
	if tx.Error != nil {
		return "", "", tx.Error
	}
	return dir, name, nil
}
