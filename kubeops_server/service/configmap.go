package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

type configmap struct{}
type CmsResp struct {
	Total int      `json:"total"`
	Item  []cmInfo `json:"item"`
}

type cmInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Modify    *bool             `json:"modify"`
	Age       string            `json:"age"`
	Data      map[string]string `json:"data"`
}

type CreateConfig struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Data      map[string]string `json:"data"`
	//是否能进行修改的标识符
	Modify bool `json:"modify"`
}

var Configmaps configmap

func (c *configmap) toCells(cms []corev1.ConfigMap) []DataCell {
	cells := make([]DataCell, len(cms))
	for i := range cms {
		cells[i] = configmapCell(cms[i])
	}
	return cells
}

func (c *configmap) fromCells(cells []DataCell) []corev1.ConfigMap {
	cms := make([]corev1.ConfigMap, len(cells))
	for i := range cells {
		cms[i] = corev1.ConfigMap(cells[i].(configmapCell))
	}
	return cms
}

// GetCmList 列表
func (c *configmap) GetCmList(configName, Namespace string, Limit, Page, uuid int) (DP *CmsResp, err error) {
	//获取deployment 的所有清单列表
	cmList, err := K8s.Clientset[uuid].CoreV1().ConfigMaps(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the ConfigMap list,reason: " + err.Error())
		return nil, err
	}
	//组装数据
	selectData := &dataselector{
		GenericDataList: c.toCells(cmList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{configName},
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
		},
	}
	//先过滤 后排序
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//分页
	dataPage := filtered.Sort().Pagination()
	cms := c.fromCells(dataPage.GenericDataList)
	item := make([]cmInfo, 0, total)
	for _, v := range cms {
		item = append(item, cmInfo{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Modify:    v.Immutable,
			Age:       v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
			Data:      v.Data,
		})
	}

	return &CmsResp{
		Total: total,
		Item:  item,
	}, nil
}

// GetCmDetail 详情
func (c *configmap) GetCmDetail(Namespace, configName string, uuid int) (detail *corev1.ConfigMap, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].CoreV1().ConfigMaps(Namespace).Get(context.TODO(), configName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the ConfigMap " + configName + " detail,reason: " + err.Error())
		return nil, err
	}
	detail.Kind = "ConfigMap"
	detail.APIVersion = "v1"
	utils.Logger.Info("Get ConfigMap " + configName + "success")
	return detail, nil
}

// CreateCm 创建
func (c *configmap) CreateCm(data *CreateConfig, uuid int) (err error) {
	config := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Immutable: &data.Modify,
		Data:      data.Data,
	}
	_, err = K8s.Clientset[uuid].CoreV1().ConfigMaps(data.Namespace).Create(context.TODO(), config, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Create ConfigMap " + data.Name + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Create ConfigMap " + data.Name + "success")
	return nil

}

// DelCm 删除
func (c *configmap) DelCm(Namespace, configName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().ConfigMaps(Namespace).Delete(context.TODO(), configName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete ConfigMap " + configName + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Delete ConfigMap " + configName + "success")
	return nil
}

// UpdateCm  更新
func (c *configmap) UpdateCm(Namespace string, Config *corev1.ConfigMap, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().ConfigMaps(Namespace).Update(context.TODO(), Config, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update ConfigMap" + Config.Name + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Update ConfigMap " + Config.Name + "success")
	return nil
}
