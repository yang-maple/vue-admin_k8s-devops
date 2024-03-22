package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

type namespace struct{}

var Namespace namespace

type NsList struct {
	Total int              `json:"total"`
	Item  []namespacesInfo `json:"item"`
}

type namespacesInfo struct {
	Name   string                `json:"name"`
	Labels map[string]string     `json:"labels"`
	Status corev1.NamespacePhase `json:"status"`
	Age    string                `json:"age"`
}

type NamespaceDetail struct {
	Detail *corev1.Namespace `json:"detail"`
	Age    string            `json:"age"`
}

func (n *namespace) toCells(ns []corev1.Namespace) []DataCell {
	cells := make([]DataCell, len(ns))
	for i := range ns {
		cells[i] = nsCell(ns[i])
	}
	return cells
}

func (n *namespace) fromCells(cells []DataCell) []corev1.Namespace {
	ns := make([]corev1.Namespace, len(cells))
	for i := range cells {
		ns[i] = corev1.Namespace(cells[i].(nsCell))
	}
	return ns
}

// GetNsList  列表
func (n *namespace) GetNsList(Name string, Limit, Page int, uuid int) (namespacesList *NsList, err error) {
	nsList, err := K8s.Clientset[uuid].CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Namespaces list,reason: " + err.Error())
		return nil, err
	}

	//组装数据
	selectData := &dataselector{
		DataSelect: &DataSelectQuery{
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
			Filter: &FilterQuery{
				Name: Name,
			},
		},
		GenericDataList: n.toCells(nsList.Items),
	}

	//筛选
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//排序并分页
	dataPage := filtered.Sort().Pagination()
	list := make([]namespacesInfo, 0, len(nsList.Items))
	for _, v := range n.fromCells(dataPage.GenericDataList) {
		list = append(list, namespacesInfo{
			Name:   v.Name,
			Labels: v.Labels,
			Status: v.Status.Phase,
			Age:    v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return &NsList{
		Total: total,
		Item:  list,
	}, nil
}

// GetNsDetail 获取namespace 详情
func (n *namespace) GetNsDetail(NsName string, uuid int) (*NamespaceDetail, error) {
	//获取deploy
	details, err := K8s.Clientset[uuid].CoreV1().Namespaces().Get(context.TODO(), NsName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Namespaces " + NsName + " detail,reason: " + err.Error())
		return nil, err
	}
	details.Kind = "Namespace"
	details.APIVersion = "v1"
	utils.Logger.Info("Get Namespaces " + NsName + "success")
	return &NamespaceDetail{
		Detail: details,
		Age:    utils.GetAge(details.CreationTimestamp.Unix()),
	}, nil
}

// DelNs 删除
func (n *namespace) DelNs(NsName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().Namespaces().Delete(context.TODO(), NsName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the Namespaces " + NsName + " detail,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete Namespaces " + NsName + "success")
	return nil
}

// CreateNs 创建
func (n *namespace) CreateNs(NsName string, uuid int) (err error) {
	namespaceConfig := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: NsName,
		},
	}
	_, err = K8s.Clientset[uuid].CoreV1().Namespaces().Create(context.TODO(), namespaceConfig, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Create ConfigMap " + NsName + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Create Namespaces " + NsName + "success")
	return nil
}

// UpdateNs  更新
func (n *namespace) UpdateNs(ns *corev1.Namespace, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().Namespaces().Update(context.TODO(), ns, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update the Namespaces " + ns.Name + " detail,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update Namespaces " + ns.Name + "success")
	return nil
}
