package service

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

// 定义空结构体
type statefulSet struct{}

// StatefulSet  全局变量供外部调用
var StatefulSet statefulSet

// StsResp 定义返回的结构体
type StsResp struct {
	Total int            `json:"total"`
	Item  []statefulInfo `json:"item"`
}

type statefulInfo struct {
	Name       string            `json:"name"`
	Namespaces string            `json:"namespaces"`
	Image      []string          `json:"image"`
	Labels     map[string]string `json:"labels"`
	Pods       string            `json:"pods"`
	Age        string            `json:"age"`
	Status     string            `json:"status"`
}

type CountStateful struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

// toCells 将 statefulCell 转换为 dataCell
func (s *statefulSet) toCells(statefulCell []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(statefulCell))
	for i := range cells {
		cells[i] = statefulsetCell(statefulCell[i])
	}
	return cells
}

// 将dataCell 转换为 statefulCell
func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulCell := make([]appsv1.StatefulSet, len(cells))
	for i := range statefulCell {
		statefulCell[i] = appsv1.StatefulSet(cells[i].(statefulsetCell))
	}
	return statefulCell
}

// GetStatefulList  列表
func (s *statefulSet) GetStatefulList(StsName, Namespace string, Limit, Page int, uuid int) (stsResp *StsResp, err error) {
	statefulList, err := K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the StatefulSets list,reason: " + err.Error())
		return nil, err
	}
	selectData := &dataselector{
		GenericDataList: s.toCells(statefulList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{StsName},
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
	states := s.fromCells(dataPage.GenericDataList)
	item := make([]statefulInfo, 0, total)
	for _, v := range states {
		images := make([]string, 0, len(v.Spec.Template.Spec.Containers))
		for _, im := range v.Spec.Template.Spec.Containers {
			images = append(images, im.Image)
		}
		pods, status := utils.GetStatus(v.Status.Replicas, v.Status.ReadyReplicas)
		item = append(item, statefulInfo{
			Name:       v.Name,
			Namespaces: v.Namespace,
			Image:      images,
			Labels:     v.Labels,
			Pods:       pods,
			Age:        v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
			Status:     status,
		})
	}
	return &StsResp{
		Total: total,
		Item:  item,
	}, nil

}

// GetStatefulDetail 详情
func (s *statefulSet) GetStatefulDetail(Namespace, StsName string, uuid int) (detail *appsv1.StatefulSet, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Get(context.TODO(), StsName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the StatefulSets " + StsName + " detail,reason: " + err.Error())
		return nil, err
	}
	detail.Kind = "StatefulSet"
	detail.APIVersion = "apps/v1"
	utils.Logger.Info("Get StatefulSets " + StsName + "success")
	return detail, nil
}

// DelStateful 删除
func (s *statefulSet) DelStateful(Namespace, StsName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Delete(context.TODO(), StsName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the StatefulSets " + StsName + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete StatefulSets " + StsName + " success")
	return nil
}

// UpdateDelStateful 更新
func (s *statefulSet) UpdateDelStateful(Namespace string, Sts *appsv1.StatefulSet, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Update(context.TODO(), Sts, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the StatefulSets " + Sts.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete StatefulSets " + Sts.Name + " success")
	return nil
}

// ModifyStatefulReplicas 修改StatefulSets 副本数
func (s *statefulSet) ModifyStatefulReplicas(Namespace, StsName string, Replicas *int32, uuid int) (err error) {
	stateful, err := K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Get(context.TODO(), StsName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to modify the number of copies of " + StsName + " ,reason: " + err.Error())
		return err
	}
	stateful.Spec.Replicas = Replicas
	_, err = K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).Update(context.TODO(), stateful, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to modify the number of copies of " + StsName + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Changed the number of copies of " + StsName + " successfully")
	return nil
}

// StatefulCount  统计stateful ready数量
func (s *statefulSet) StatefulCount(Namespace string, uuid int) (*CountStateful, error) {
	statefulList, err := K8s.Clientset[uuid].AppsV1().StatefulSets(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to count the number of stators", err.Error())
		return nil, err
	}
	count := 0
	for _, v := range statefulList.Items {
		if v.Status.ReadyReplicas == *v.Spec.Replicas {
			count++
		}
	}
	return &CountStateful{
		Total:    len(statefulList.Items),
		Ready:    count,
		NotReady: len(statefulList.Items) - count,
	}, nil
}
