package service

import (
	"github.com/gin-gonic/gin"
	"kubeops/dao"
	"kubeops/utils"
	"mime/multipart"
	"os"
	"path"
)

var Cluster cluster

type cluster struct{}
type clusterList struct {
	Item  []clusterInfo `json:"item"`
	Total int           `json:"total"`
}

type clusterInfo struct {
	Id          uint   `json:"id"`
	ClusterName string `json:"cluster_name"`
	FileName    string `json:"file_name"`
	Type        string `json:"type"`
	Status      bool   `json:"status"`
	CreateTime  string `json:"create_time"`
}

// Create 创建集群
func (clt *cluster) Create(dir, clusterName, clusterType string, uuid int, file *multipart.FileHeader, c *gin.Context) error {
	//判断目录是不是存在,如果存在则不创建
	if !utils.DirExists(dir) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			utils.Logger.Error("Failed to create a directory,reason:" + err.Error())
			return err
		}
	}
	//保留原始文件名
	oldFilename := file.Filename
	//重新命名文件
	file.Filename = clusterName + "_" + clusterType + ".conf"
	//保存文件
	dst := path.Join(dir, file.Filename)
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		utils.Logger.Error("Failed to save file,reason:" + err.Error())
		return err
	}
	//数据库存储
	err = dao.Cluster.AddCluster(clusterName, oldFilename, clusterType, dir, uuid)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}
	utils.Logger.Info("The cluster " + clusterName + " was successfully added")
	return nil
}

// List 获取集群列表
func (clt *cluster) List(uuid int) (item *clusterList, err error) {
	//获取用户集群列表
	value, err := dao.Cluster.ListCluster(uuid)
	if err != nil {
		utils.Logger.Error("Failed to get the Cluster,reason: " + err.Error())
		return nil, err
	}
	clusters := make([]clusterInfo, 0, len(*value))
	for _, v := range *value {
		cluster := clusterInfo{
			ClusterName: v.ClusterName,
			FileName:    v.FileName,
			Type:        v.Type,
			Status:      v.Status,
			CreateTime:  v.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		clusters = append(clusters, cluster)
	}
	//返回集群列表
	return &clusterList{
		Item:  clusters,
		Total: len(clusters),
	}, nil
}

// GetClusterDetail 获取集群信息
func (clt *cluster) GetClusterDetail(clusterName string, uuid int) (item *clusterInfo, err error) {
	//获取集群信息
	value, err := dao.Cluster.GetClusterDetail(clusterName, uuid)
	if err != nil {
		utils.Logger.Error("Failed to get the cluster detail,reason: " + err.Error())
		return nil, err
	}
	//返回集群信息
	return &clusterInfo{
		Id:          value.Id,
		ClusterName: value.ClusterName,
		Type:        value.Type,
	}, err
}

// Change 更换集群信息
func (clt *cluster) Change(clusterName string, uuid int) error {
	//找到dir - 初始化 - 更改
	// 查询dir
	dir, err := dao.Cluster.GetClusterDir(clusterName, uuid)
	if err != nil {
		utils.Logger.Error("Failed to change the cluster,reason: " + err.Error())
		return err
	}
	//更新k8s client set
	K8s.ConfigDir[uuid] = &dir
	err = K8s.Init(uuid)
	if err != nil {
		utils.Logger.Error("Failed to initialize the cluster,reason: " + err.Error())
		return err
	}
	//更改sql 中集群的状态
	err = dao.Cluster.ChangeCluster(clusterName, uuid)
	if err != nil {
		utils.Logger.Error("Failed to change the cluster status,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("The cluster " + clusterName + " was successfully changed")
	return nil
}

// Delete 删除集群
func (clt *cluster) Delete(name string, userid int) error {
	//获取 dir -- 删除本地 -- 删除库
	//获取集群路径
	dir, err := dao.Cluster.GetClusterDir(name, userid)
	if err != nil {
		utils.Logger.Error("Failed to delete the cluster ,reason: " + err.Error())
		return err
	}
	if dir != "" {
		//删除本地文件
		err = os.RemoveAll(dir)
		if err != nil {
			utils.Logger.Error("Failed to delete the file ,reason: " + err.Error())
			return err
		}
	}
	err = dao.Cluster.DeleteCluster(name, userid)
	if err != nil {
		utils.Logger.Error("Failed to delete the file ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("The cluster " + name + " was successfully deleted")
	return nil
}

// UpdateCluster 更新集群信息
func (clt *cluster) UpdateCluster(Id uint, cluster, clusterType string, uuid int) error {
	//更改数据库文件名称
	dir, oldFileName, err := dao.Cluster.UpdateCluster(Id, cluster, clusterType, uuid)
	if err != nil {
		utils.Logger.Error("Failed to update the data cluster information ,reason: " + err.Error())
		return err
	}
	newFileName := dir + "/" + cluster + "_" + clusterType + ".conf"
	//更改存储文件名称
	err = os.Rename(oldFileName, newFileName)
	if err != nil {
		utils.Logger.Error("Failed to update the data cluster name ,reason: " + err.Error())
		return err
	}
	//更改配置文件路径
	K8s.ConfigDir[uuid] = &newFileName
	utils.Logger.Info("The cluster " + cluster + " was successfully updated")
	return nil
}
