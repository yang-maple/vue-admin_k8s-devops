package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

var StorageClass storageClass

type storageClass struct {
}

type storageClassResp struct {
	Total int                `json:"total"`
	Item  []storageClassInfo `json:"item"`
}

type storageClassInfo struct {
	Name          string                                `json:"name"`
	Provisioner   string                                `json:"provisioner"`
	Parameters    map[string]string                     `json:"parameters"`
	ReclaimPolicy *corev1.PersistentVolumeReclaimPolicy `json:"reclaim_policy"`
	Age           string                                `json:"age"`
}

func (sc *storageClass) toCells(storage []storagev1.StorageClass) []DataCell {
	cells := make([]DataCell, len(storage))
	for i := range storage {
		cells[i] = storageClassCell(storage[i])
	}
	return cells
}

func (sc *storageClass) fromCells(cells []DataCell) []storagev1.StorageClass {
	storages := make([]storagev1.StorageClass, len(cells))
	for i := range cells {
		storages[i] = storagev1.StorageClass(cells[i].(storageClassCell))
	}
	return storages
}

// GetStorageClassList 获取存储类列表
func (sc *storageClass) GetStorageClassList(storageName string, Limit, Page int, uuid int) (*storageClassResp, error) {
	storageClasses, err := K8s.Clientset[uuid].StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the StorageClasses list,reason: " + err.Error())
		return nil, err
	}

	//组装数据
	selectData := &dataselector{
		GenericDataList: sc.toCells(storageClasses.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{storageName},
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
		},
	}

	//筛选数据
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//排序并分页
	dataPage := filtered.Sort().Pagination()
	storages := sc.fromCells(dataPage.GenericDataList)
	item := make([]storageClassInfo, 0, len(storageClasses.Items))
	for _, storageClass := range storages {
		item = append(item, storageClassInfo{
			Name:          storageClass.Name,
			Provisioner:   storageClass.Provisioner,
			Parameters:    storageClass.Parameters,
			ReclaimPolicy: storageClass.ReclaimPolicy,
			Age:           storageClass.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return &storageClassResp{
		Total: total,
		Item:  item,
	}, nil
}

// GetStorageClassDetail 获取存储类详情
func (sc *storageClass) GetStorageClassDetail(storageClassName string, uuid int) (*storagev1.StorageClass, error) {
	storageClass, err := K8s.Clientset[uuid].StorageV1().StorageClasses().Get(context.TODO(), storageClassName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the StorageClasses " + storageClassName + " detail,reason: " + err.Error())
		return nil, err
	}
	storageClass.APIVersion = "storage.k8s.io/v1"
	storageClass.Kind = "StorageClass"
	utils.Logger.Info("Get StorageClasses " + storageClassName + "success")
	return storageClass, nil
}

// DelStorageClass 删除存储类
func (sc *storageClass) DelStorageClass(storageClassName string, uuid int) error {
	err := K8s.Clientset[uuid].StorageV1().StorageClasses().Delete(context.TODO(), storageClassName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the StorageClasses " + storageClassName + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete StorageClasses " + storageClassName + " success")
	return nil
}

// UpdateStorageClass 更新存储类
func (sc *storageClass) UpdateStorageClass(storageClass *storagev1.StorageClass, uuid int) error {
	_, err := K8s.Clientset[uuid].StorageV1().StorageClasses().Update(context.TODO(), storageClass, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update the StorageClasses " + storageClass.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update StorageClasses " + storageClass.Name + " success")
	return nil
}
