package service

import (
	"context"
	networkv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

var IngressClass ingressClass

type ingressClass struct {
}

type IngressClassInfo struct {
	Total int                `json:"total"`
	Item  []ingressClassInfo `json:"item"`
}

type ingressClassInfo struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
	Age        string `json:"age"`
}

func (ic *ingressClass) toCells(ingClass []networkv1.IngressClass) []DataCell {
	cells := make([]DataCell, len(ingClass))
	for i := range ingClass {
		cells[i] = ingressClassCell(ingClass[i])
	}
	return cells
}

func (ic *ingressClass) fromCells(cells []DataCell) []networkv1.IngressClass {
	ingClass := make([]networkv1.IngressClass, len(cells))
	for i := range cells {
		ingClass[i] = networkv1.IngressClass(cells[i].(ingressClassCell))
	}
	return ingClass
}

// GetIngressClass 获取IngressClass list
func (ic *ingressClass) GetIngressClass(Name string, Limit, Page int, uuid int) (ingressClass *IngressClassInfo, err error) {
	ingressClasses, err := K8s.Clientset[uuid].NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the IngressClasses list,reason: " + err.Error())
		return nil, err
	}
	//组装数据
	selectData := &dataselector{
		GenericDataList: ic.toCells(ingressClasses.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{
				Name: Name,
			},
			Paginate: &PaginateQuery{
				page:  Page,
				limit: Limit,
			},
		},
	}
	//筛选
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//排序分页
	dataPage := filtered.Sort().Pagination()
	item := make([]ingressClassInfo, 0, len(ingressClasses.Items))
	for _, v := range ic.fromCells(dataPage.GenericDataList) {
		item = append(item, ingressClassInfo{
			Name:       v.Name,
			Controller: v.Spec.Controller,
			Age:        v.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return &IngressClassInfo{
		Total: total,
		Item:  item,
	}, nil
}

// GetIngressClassDetail 获取IngressClass详情
func (ic *ingressClass) GetIngressClassDetail(Name string, uuid int) (ingressClass *networkv1.IngressClass, err error) {
	ingressClass, err = K8s.Clientset[uuid].NetworkingV1().IngressClasses().Get(context.TODO(), Name, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the IngressClasses " + Name + " detail,reason: " + err.Error())
		return nil, err
	}
	ingressClass.Kind = "IngressClass"
	ingressClass.APIVersion = "networking.k8s.io/v1"
	utils.Logger.Info("Get IngressClasses " + Name + "success")
	return ingressClass, nil
}

// DelIngressClass 删除ingressClass
func (ic *ingressClass) DelIngressClass(Name string, uuid int) (err error) {
	err = K8s.Clientset[uuid].NetworkingV1().IngressClasses().Delete(context.TODO(), Name, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the IngressClasses " + Name + " detail,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete IngressClasses " + Name + "success")
	return nil
}

// UpdateIngressClass 更新IngressClass
func (ic *ingressClass) UpdateIngressClass(ingressClass *networkv1.IngressClass, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].NetworkingV1().IngressClasses().Update(context.TODO(), ingressClass, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update the IngressClasses " + ingressClass.Name + " detail,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update IngressClasses " + ingressClass.Name + "success")
	return nil
}
