package service

//列表
// 获取node 详情
import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
	"strconv"
)

type node struct{}

var Node node

type NodeList struct {
	Total     int        `json:"total"`
	Item      []nodeInfo `json:"Item"`
	Resources *Resource  `json:"resource"`
}

type Resource struct {
	CpuRequest    map[string]float64 `json:"cpu_request"`
	CpuLimit      map[string]float64 `json:"cpu_limit"`
	MemoryRequest map[string]int64   `json:"memory_request"`
	MemoryLimit   map[string]int64   `json:"memory_limit"`
}

type nodeInfo struct {
	Name           string            `json:"name"`
	Labels         map[string]string `json:"labels"`
	Status         string            `json:"status"`
	Unschedulable  bool              `json:"unschedulable"`
	Taints         []corev1.Taint    `json:"taints"`
	NodeIp         string            `json:"nodeIp"`
	CpuTotal       float64           `json:"cpu_total"`
	MemoryTotal    int64             `json:"memory_total"`
	Pods           int               `json:"pods"`
	CreateTime     string            `json:"create_time"`
	KubeletVersion string            `json:"kubelet_version"`
}

type NodeDetail struct {
	Detail          *corev1.Node `json:"detail"`
	Pods            []poddetail  `json:"pods"`
	MemoryAllocator string       `json:"memory_allocator"`
	AGE             string       `json:"age"`
	Total           int          `json:"total"`
}

type poddetail struct {
	Name    string
	Image   []string
	Labels  map[string]string
	Status  corev1.PodPhase
	Restart int32
	PodAge  string
}

const (
	// MBMemory MB 内存 字节转MB
	MBMemory = 1048576
	// SmallCore 核数 大核心数 转小核心数
	SmallCore = 1000
)

// GetNodeList  列表
func (n *node) GetNodeList(uuid int) (nodeList *NodeList, err error) {
	nodeLists, err := K8s.Clientset[uuid].CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Nodes list,reason: " + err.Error())
		return nil, err
	}
	item := make([]nodeInfo, 0, len(nodeLists.Items))
	for _, v := range nodeLists.Items {
		status := "Ready"
		if v.Status.Conditions[len(v.Status.Conditions)-1].Status != "True" {
			status = "NotReady"
		}
		item = append(item, nodeInfo{
			Name:           v.Name,
			Labels:         v.Labels,
			Status:         status,
			Unschedulable:  v.Spec.Unschedulable,
			Taints:         v.Spec.Taints,
			NodeIp:         v.Status.Addresses[0].Address,
			CpuTotal:       float62(v.Status.Allocatable.Cpu().AsApproximateFloat64()),
			MemoryTotal:    v.Status.Allocatable.Memory().Value() / MBMemory,
			Pods:           len(*GetNodePods(v.Name, uuid)),
			CreateTime:     v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
			KubeletVersion: v.Status.NodeInfo.KubeletVersion,
		})
	}
	return &NodeList{
		Total:     len(nodeLists.Items),
		Item:      item,
		Resources: Node.GetNodeResource(uuid),
	}, nil
}

// GetNodeDetail  获取node 详情
func (n *node) GetNodeDetail(NodeName string, uuid int) (details *NodeDetail, err error) {
	//获取deploy
	detail, err := K8s.Clientset[uuid].CoreV1().Nodes().Get(context.TODO(), NodeName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Nodes " + NodeName + " detail,reason: " + err.Error())
		return nil, err
	}
	pods := GetNodePods(NodeName, uuid)
	memory, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", float64(detail.Status.Allocatable.Name(corev1.ResourceMemory, "").Value())/float64(1024*1024*1024)), 64)
	utils.Logger.Info("Get Nodes " + NodeName + "success")
	return &NodeDetail{
		Detail:          detail,
		Pods:            *pods,
		MemoryAllocator: strconv.FormatFloat(memory, 'f', -1, 64) + "Gi",
		AGE:             utils.GetAge(detail.CreationTimestamp.Unix()),
		Total:           len(*pods),
	}, nil
}

// GetNodePods 统计Pod 数量
func GetNodePods(NodeName string, uuid int) (detail *[]poddetail) {
	//获取pod list
	podList, err := K8s.Clientset[uuid].CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to count the number of pods" + err.Error())
		return nil
	}
	nodePods := make([]poddetail, 0)
	for _, pods := range podList.Items {
		if pods.Spec.NodeName == NodeName {
			containers := make([]string, 0)
			for _, container := range pods.Spec.Containers {
				containers = append(containers, container.Image)
			}
			if pods.Status.Phase != corev1.PodSucceeded {
				nodePods = append(nodePods, poddetail{
					Name:   pods.Name,
					Image:  containers,
					Labels: pods.Labels,
					Status: pods.Status.Phase,
					//Restart: pods.Status.ContainerStatuses[0].RestartCount,
					PodAge: utils.GetAge(pods.CreationTimestamp.Unix()),
				})
			}

		}
	}
	return &nodePods
}

// SetNodeSchedule 设置节点是否可调度
func (n *node) SetNodeSchedule(name string, status bool, uuid int) (err error) {
	node, err := K8s.Clientset[uuid].CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to get node status" + err.Error())
		return err
	}
	node.Spec.Unschedulable = status
	_, err = K8s.Clientset[uuid].CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to set node " + name + " status Unschedulable = " + strconv.FormatBool(status) + " :" + err.Error())
		return err
	}
	utils.Logger.Info("Succeed to set node " + name + " status Unschedulable = " + strconv.FormatBool(status))
	return nil
}

//// EmptyNode 排空节点操作
//func (n *node) EmptyNode(name string, uuid int) (err error) {
//	//获取node
//	return nil
//}

// GetNodeResource 获取节点的资源,请求，limit,使用量
func (n *node) GetNodeResource(uuid int) *Resource {
	cpuLimit := make(map[string]float64)
	cpuRequest := make(map[string]float64)
	memoryLimit := make(map[string]int64)
	memoryRequest := make(map[string]int64)
	pods, _ := K8s.Clientset[uuid].CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	for _, pod := range pods.Items {
		if pod.Spec.NodeName != "" {
			for _, container := range pod.Spec.Containers {
				if _, ok := cpuLimit[pod.Spec.NodeName]; ok {
					cpuLimit[pod.Spec.NodeName] += float62(container.Resources.Limits.Cpu().AsApproximateFloat64())
					cpuRequest[pod.Spec.NodeName] += float62(container.Resources.Requests.Cpu().AsApproximateFloat64())
					memoryLimit[pod.Spec.NodeName] += container.Resources.Limits.Memory().Value() / MBMemory
					memoryRequest[pod.Spec.NodeName] += container.Resources.Requests.Memory().Value() / MBMemory
					continue
				} else {
					cpuLimit[pod.Spec.NodeName] += float62(container.Resources.Limits.Cpu().AsApproximateFloat64())
					cpuRequest[pod.Spec.NodeName] += float62(container.Resources.Requests.Cpu().AsApproximateFloat64())
					memoryLimit[pod.Spec.NodeName] += container.Resources.Limits.Memory().Value() / MBMemory
					memoryRequest[pod.Spec.NodeName] += container.Resources.Requests.Memory().Value() / MBMemory
				}
			}
		}
	}

	//

	//百分比
	return &Resource{
		CpuRequest:    cpuRequest,
		CpuLimit:      cpuLimit,
		MemoryLimit:   memoryLimit,
		MemoryRequest: memoryRequest,
	}
}

// 保留2位小数
func float62(f float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f*SmallCore), 64)
	return value
}
