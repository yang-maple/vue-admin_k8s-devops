package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Homepage homepage

type homepage struct {
}

type homepageInfo struct {
	ClusterInfo  []clusterInfo     `json:"cluster_info"`
	NodeInfo     []nodeInfo        `json:"node_info"`
	NodeResource *Resource         `json:"node_resource"`
	DeployInfo   *CountDeploy      `json:"deploy_info"`
	StatefulInfo *CountStateful    `json:"stateful_info"`
	DaemonInfo   *CountDaemon      `json:"daemon_info"`
	PodInfo      *CountPodReady    `json:"pod_info"`
	PodTotal     map[string]nsInfo `json:"pod_total"`
	EventTotal   EventResp         `json:"event_total"`
}

type nsInfo struct {
	PodTotal      int `json:"pod_total"`
	DeployTotal   int `json:"deploy_total"`
	StatefulTotal int `json:"stateful_total"`
	DaemonTotal   int `json:"daemon_total"`
	JobTotal      int `json:"job_total"`
}

type homeJob struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

// GetHomepage 获取首页信息
func (h *homepage) GetHomepage(uuid int) (info *homepageInfo, err error) {
	cInfo, err := Cluster.List(uuid)
	if err != nil {
		return nil, err
	}
	nInfo, err := Node.GetNodeList(uuid)
	if err != nil {
		return nil, err
	}
	dInfo, err := Deployment.DeployCount("", uuid)
	if err != nil {
		return nil, err
	}
	sInfo, err := StatefulSet.StatefulCount("", uuid)
	if err != nil {
		return nil, err
	}
	daeInfo, err := DaemonSet.DaemonCount("", uuid)
	if err != nil {
		return nil, err
	}
	podsInfo, err := Pod.PodCount("", uuid)
	podCountPerNamespace, err := getlist(uuid)
	if err != nil {
		return nil, err
	}
	eventTotal, err := Event.ListEvent(uuid, "Normal", "")
	if err != nil {
		return nil, err
	}
	return &homepageInfo{
		ClusterInfo:  cInfo.Item,
		NodeInfo:     nInfo.Item,
		NodeResource: nInfo.Resources,
		DeployInfo:   dInfo,
		StatefulInfo: sInfo,
		DaemonInfo:   daeInfo,
		PodInfo:      podsInfo,
		PodTotal:     podCountPerNamespace,
		EventTotal:   *eventTotal,
	}, nil
}

// 获取ns每个资源
func getlist(uuid int) (map[string]nsInfo, error) {
	// 统计每个ns的pod数量
	podCountPerNamespace := make(map[string]nsInfo)
	pod, err := K8s.Clientset[uuid].CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	//检查podCountPerNamespace中是否已经包含了特定命名空间的条目。
	for _, v := range pod.Items {
		currentNsInfo, ok := podCountPerNamespace[v.Namespace]
		if ok {
			//如果已经存在，则增加PodTotal计数，并保留其他类型的计数不变。
			currentNsInfo.PodTotal++
			podCountPerNamespace[v.Namespace] = currentNsInfo
		} else {
			//如果不存在，则创建一个新的nsInfo对象，并将PodTotal设置为1，并保留其他类型的计数不变
			currentNsInfo = nsInfo{
				PodTotal: 1,
			}
			podCountPerNamespace[v.Namespace] = currentNsInfo
		}
	}
	// 统计每个ns的deployment数量
	deploy, err := K8s.Clientset[uuid].AppsV1().Deployments(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range deploy.Items {
		currentNsInfo, ok := podCountPerNamespace[v.Namespace]
		if ok {
			currentNsInfo.DeployTotal++
			podCountPerNamespace[v.Namespace] = currentNsInfo
		} else {
			currentNsInfo = nsInfo{
				DeployTotal: 1,
			}
			podCountPerNamespace[v.Namespace] = currentNsInfo
		}
	}
	// 统计每个ns的stateful数量
	stateful, err := K8s.Clientset[uuid].AppsV1().StatefulSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range stateful.Items {
		currentNsInfo, ok := podCountPerNamespace[v.Namespace]
		if ok {
			currentNsInfo.StatefulTotal++
			podCountPerNamespace[v.Namespace] = currentNsInfo
		} else {
			currentNsInfo = nsInfo{
				StatefulTotal: 1,
			}
			podCountPerNamespace[v.Namespace] = currentNsInfo
		}
	}
	// 统计每个ns的daemon数量
	daemon, err := K8s.Clientset[uuid].AppsV1().DaemonSets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range daemon.Items {
		currentNsInfo, ok := podCountPerNamespace[v.Namespace]
		if ok {
			currentNsInfo.DaemonTotal++
			podCountPerNamespace[v.Namespace] = currentNsInfo
		} else {
			currentNsInfo = nsInfo{
				DaemonTotal: 1,
			}
			podCountPerNamespace[v.Namespace] = currentNsInfo
		}
	}
	// 统计每个ns的job数量
	job, err := K8s.Clientset[uuid].BatchV1().Jobs(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, v := range job.Items {
		currentNsInfo, ok := podCountPerNamespace[v.Namespace]
		if ok {
			currentNsInfo.JobTotal++
			podCountPerNamespace[v.Namespace] = currentNsInfo
		} else {
			currentNsInfo = nsInfo{
				JobTotal: 1,
			}
			podCountPerNamespace[v.Namespace] = currentNsInfo
		}
	}
	// 返回结果
	return podCountPerNamespace, err
}

func getevent(uuid int) (info string, err error) {
	K8s.Clientset[uuid].CoreV1().Events(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	return "", nil
}
