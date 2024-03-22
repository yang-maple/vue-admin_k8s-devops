package service

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubeops/utils"
)

type service struct{}

var Services service

type svcResp struct {
	Total int       `json:"total"`
	Item  []svcInfo `json:"item"`
}

type svcInfo struct {
	Name       string               `json:"name"`
	Namespace  string               `json:"namespace"`
	Labels     map[string]string    `json:"labels"`
	Type       corev1.ServiceType   `json:"type"`
	ClusterIp  string               `json:"cluster_ip"`
	ExternalIp []string             `json:"external_ip"`
	Port       []corev1.ServicePort `json:"port"`
	Age        string               `json:"age"`
}

type CreateService struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace"`
	Labels       map[string]string `json:"labels"`
	Type         string            `json:"type"`
	ServicePorts Port              `json:"service_ports"`
}
type Port struct {
	PortName   string          `json:"port_name"`
	Port       int32           `json:"port"`
	Protocol   corev1.Protocol `json:"protocol"`
	TargetPort int32           `json:"target_port"`
	NodePort   int32           `json:"node_port"`
}

type svcDetail struct {
	Detail *corev1.Service `json:"detail"`
	Age    string          `json:"age"`
}

// 数据类型转换 cells
func (s *service) toCells(svc []corev1.Service) []DataCell {
	cells := make([]DataCell, len(svc))
	for i := range svc {
		cells[i] = serviceCell(svc[i])
	}
	return cells
}

func (s *service) fromCells(cells []DataCell) []corev1.Service {
	svc := make([]corev1.Service, len(cells))
	for i := range cells {
		svc[i] = corev1.Service(cells[i].(serviceCell))
	}
	return svc
}

// GetSvcList 获取services列表
func (s *service) GetSvcList(svcName, Namespace string, Limit, Page int, uuid int) (DP *svcResp, err error) {
	//获取deployment 的所有清单列表
	svcList, err := K8s.Clientset[uuid].CoreV1().Services(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Services list,reason: " + err.Error())
		return nil, err
	}
	//组装数据
	selectData := &dataselector{
		GenericDataList: s.toCells(svcList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{svcName},
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
	svc := s.fromCells(dataPage.GenericDataList)
	item := make([]svcInfo, 0, total)
	for _, v := range svc {
		externalIp := make([]string, 0)
		for _, value := range v.Spec.ExternalIPs {
			externalIp = append(externalIp, value)

		}
		item = append(item, svcInfo{
			Name:       v.Name,
			Namespace:  v.Namespace,
			Labels:     v.Labels,
			Type:       v.Spec.Type,
			ClusterIp:  v.Spec.ClusterIP,
			ExternalIp: externalIp,
			Port:       v.Spec.Ports,
			Age:        v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return &svcResp{
		Total: total,
		Item:  item,
	}, nil
}

// GetSvcDetail 获取 services 详情
func (s *service) GetSvcDetail(Namespace, svcName string, uuid int) (*svcDetail, error) {
	//获取deploy
	detail, err := K8s.Clientset[uuid].CoreV1().Services(Namespace).Get(context.TODO(), svcName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Services " + svcName + " detail,reason: " + err.Error())
		return nil, err
	}
	detail.Kind = "Service"
	detail.APIVersion = "v1"
	utils.Logger.Info("Get Services " + svcName + "success")
	return &svcDetail{
		Detail: detail,
		Age:    utils.GetAge(detail.CreationTimestamp.Unix()),
	}, nil
}

// CreateSvc 创建 services
func (s *service) CreateSvc(data *CreateService, uuid int) (err error) {
	createSvc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: data.Labels,
			Type:     "",
		},
		Status: corev1.ServiceStatus{},
	}
	specPort := corev1.ServicePort{
		Name:     data.ServicePorts.PortName,
		Protocol: data.ServicePorts.Protocol,
		Port:     data.ServicePorts.Port,
		TargetPort: intstr.IntOrString{
			Type:   0,
			IntVal: data.ServicePorts.TargetPort,
		},
		NodePort: 0,
	}
	if data.Type == "NodePort" && data.ServicePorts.NodePort != 0 {
		createSvc.Spec.Type = "NodePort"
		specPort.NodePort = data.ServicePorts.NodePort
	}
	createSvc.Spec.Ports = append(createSvc.Spec.Ports, specPort)

	_, err = K8s.Clientset[uuid].CoreV1().Services(data.Namespace).Create(context.TODO(), createSvc, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Create the Services " + data.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Create Services " + data.Name + " success")
	return nil

}

// DelSvc 删除 services
func (s *service) DelSvc(Namespace, svcName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().Services(Namespace).Delete(context.TODO(), svcName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the Services " + svcName + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete Services " + svcName + " success")
	return nil
}

// UpdateSvc 更新 services
func (s *service) UpdateSvc(Namespace string, svc *corev1.Service, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().Services(Namespace).Update(context.TODO(), svc, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update the Services " + svc.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update Services " + svc.Name + " success")
	return nil
}
