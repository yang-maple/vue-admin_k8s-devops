package service

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kubeops/utils"
	"strconv"
	"time"
)

// 定义一个deployment 结构体
type deployment struct {
}

// Deployment 定义一个 deployment 类型的全局变量 Deployment
var Deployment deployment

// apps1.deployment --> Data cell
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	//数据转换，将pods类型转化为data cells 类型
	for i := range deployments {
		cells[i] = deploymentCell(deployments[i])
	}
	return cells
}

// Data cell转换成 apps1.deployment 类型
func (d *deployment) fromCells(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

// DeployResp 定义获取清单的列表
type DeployResp struct {
	Total int          `json:"total"`
	Item  []deployInfo `json:"item"`
}

type deployInfo struct {
	Name       string            `json:"name"`
	Namespaces string            `json:"namespaces"`
	Image      []string          `json:"image"`
	Labels     map[string]string `json:"labels"`
	Pods       string            `json:"pods"`
	Age        string            `json:"age"`
	Status     string            `json:"status"`
}

// DeploymentCreate 定义DeploymentCreate 结构体用于创建 deployment需要的参数属性的定义
type DeploymentCreate struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Replicas    int32             `json:"replicas"`
	Labels      map[string]string `json:"labels"`
	Container   container         `json:"container"`
	HealthCheck bool              `json:"health_check"`
	HealthPath  string            `json:"health_path"`
}

type container struct {
	ContainerName string        `json:"container_name"`
	Image         string        `json:"image"`
	Cpu           string        `json:"cpu"`
	Memory        string        `json:"memory"`
	Containerport containerPort `json:"container_port"`
}

type containerPort struct {
	PortName      string          `json:"port_name"`
	ContainerPort int32           `json:"container_port"`
	Protocol      corev1.Protocol `json:"protocol"`
}

// DeployDp  定义返回的数据类型
type DeployDp struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

type CountDeploy struct {
	Total    int `json:"total"`
	Ready    int `json:"ready"`
	NotReady int `json:"not_ready"`
}

// GetDeploymentList 获取deployment 列表
func (d *deployment) GetDeploymentList(DeployName, Namespace string, Limit, Page int, uuid int) (DP *DeployResp, err error) {
	//获取deployment 的所有清单列表
	deployList, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Deployments list,reason: " + err.Error())
		return nil, err
	}
	//组装数据
	selectData := &dataselector{
		GenericDataList: d.toCells(deployList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{DeployName},
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
	deploys := d.fromCells(dataPage.GenericDataList)
	item := make([]deployInfo, 0, len(deploys))
	for _, v := range deploys {
		images := make([]string, 0, len(v.Spec.Template.Spec.Containers))
		for _, im := range v.Spec.Template.Spec.Containers {
			images = append(images, im.Image)
		}
		pods, status := utils.GetStatus(v.Status.Replicas, v.Status.ReadyReplicas)
		item = append(item, deployInfo{
			Name:       v.Name,
			Namespaces: v.Namespace,
			Image:      images,
			Labels:     v.Labels,
			Pods:       pods,
			Age:        v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
			Status:     status,
		})
	}

	return &DeployResp{
		Total: total,
		Item:  item,
	}, nil
}

// GetDeployDetail 获取deployment 详情
func (d *deployment) GetDeployDetail(Namespace, DeployName string, uuid int) (detail *appsv1.Deployment, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Deployments " + DeployName + " detail,reason: " + err.Error())
		return nil, err
	}
	detail.Kind = "Deployment"
	detail.APIVersion = "apps/v1"
	utils.Logger.Info("Get Deployments " + DeployName + "success")
	return detail, nil
}

// ModifyDeployReplicas 修改deployment 副本数
func (d *deployment) ModifyDeployReplicas(Namespace, DeployName string, Replicas *int32, uuid int) (err error) {
	deploy, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to update the number of replicas with " + DeployName + ",reason: " + err.Error())
		return err
	}
	deploy.Spec.Replicas = Replicas
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to update the number of replicas with " + DeployName + ",reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update the number of replicas with " + DeployName + " success")
	return nil
}

// CreateDeploy 创建 deployment实例
func (d *deployment) CreateDeploy(data *DeploymentCreate, uuid int) (err error) {

	containers := make([]corev1.Container, 0, 1)
	containersPort := make([]corev1.ContainerPort, 0, 1)
	containersPort = append(containersPort, corev1.ContainerPort{
		Name:          data.Container.Containerport.PortName,
		ContainerPort: data.Container.Containerport.ContainerPort,
		Protocol:      data.Container.Containerport.Protocol,
	})
	containers = append(containers, corev1.Container{
		Name:  data.Container.ContainerName,
		Image: data.Container.Image,
		Ports: containersPort,
		Resources: corev1.ResourceRequirements{
			Limits: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(data.Container.Cpu),
				corev1.ResourceMemory: resource.MustParse(data.Container.Memory),
			},
			Requests: map[corev1.ResourceName]resource.Quantity{
				corev1.ResourceCPU:    resource.MustParse(data.Container.Cpu),
				corev1.ResourceMemory: resource.MustParse(data.Container.Memory),
			},
		},
	})
	//默认添加app标签
	labels := map[string]string{
		"app": data.Name,
	}
	//如果用户自定义了标签，则添加标签
	if data.Labels != nil {
		for k, v := range data.Labels {
			labels[k] = v
		}
	}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers:     containers,
				},
			},
		},
	}
	// 判断是否打开健康检测
	if data.HealthCheck {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			//设置第一个容器的readiness probe
			// 若存在多个容器 使用for 进行定义
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						// intstr.IntOrString 的作用是端口可以定义为整型，也可以定义为字符
						// type=0 表示该结构体实例内的数据为整型，转json 时只需要使用 IntVal 的数据
						// type=1 表示该结构体实例内的数据为字符串，转json 时只需要使用 StrVal 的数据
						Type:   0,
						IntVal: data.Container.Containerport.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(data.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Create Deployments " + data.Name + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Create Deployments " + data.Name + " success")
	return nil

}

// DelDeploy 删除 deployment 实例
func (d *deployment) DelDeploy(Namespace, DeployName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Delete(context.TODO(), DeployName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete Deployments " + DeployName + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Delete Deployments " + DeployName + " success")
	return nil
}

// RestartDeploy 重启 deployment 实例
func (d *deployment) RestartDeploy(Namespace, DeployName string, uuid int) (err error) {
	restartTime := fmt.Sprintf(`{"spec":{"template":{"metadata":{"labels":{"restart-time":"%s"}}}}}`, time.Now().Format("2006-01-02_15-04-05"))
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Patch(context.TODO(), DeployName, types.StrategicMergePatchType, []byte(restartTime), metav1.PatchOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Restart Deployments " + DeployName + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Restart Deployments " + DeployName + " success")
	return nil
}

// UpdateDeploy 更新 deployment 实例
func (d *deployment) UpdateDeploy(Namespace string, deploy *appsv1.Deployment, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update Deployments " + deploy.Name + ",reason:" + err.Error())
		return err
	}
	utils.Logger.Info("Update Deployments " + deploy.Name + " success")
	return nil
}

// GetDeployPer 获取 每个namespace 下的deployment 实例
func (d *deployment) GetDeployPer(uuid int) (dps []DeployDp, err error) {

	namespaceList, err := K8s.Clientset[uuid].CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to get deployment instances under each namespace,reason: " + err.Error())
		return nil, err
	}

	for _, v := range namespaceList.Items {
		deployList, err := K8s.Clientset[uuid].AppsV1().Deployments(v.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			utils.Logger.Error("Failed to get deployment instances under each namespace,reason: " + err.Error())
			return nil, err
		}
		np := &DeployDp{
			Name:  v.Name,
			Total: len(deployList.Items),
		}
		dps = append(dps, *np)
	}
	return dps, nil
}

// RolloutDeploy 回滚 deployment
func (d *deployment) RolloutDeploy(Namespace, DeployName string, uuid int) (err error) {
	//获取deploy 信息
	deploy, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Get(context.TODO(), DeployName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to roll back the deployment " + DeployName + ",reason: " + err.Error())
		return err
	}
	//获取当前版本号
	revision := deploy.Annotations["deployment.kubernetes.io/revision"]
	//回滚至上一个版本
	version, _ := strconv.Atoi(revision)
	patchData := fmt.Sprintf(`{"metadata":{"annotations":{"deployment.kubernetes.io/revision":"%d"}}}`, version)
	_, err = K8s.Clientset[uuid].AppsV1().Deployments(Namespace).Patch(context.TODO(), DeployName, types.StrategicMergePatchType, []byte(patchData), metav1.PatchOptions{})
	if err != nil {
		utils.Logger.Error("Failed to roll back the deployment " + DeployName + ",reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Roll back the deployment " + DeployName + " success")
	return nil
}

// DeployCount 统计deployment 实例的状态分布
func (d *deployment) DeployCount(Namespace string, uuid int) (*CountDeploy, error) {
	deployList, err := K8s.Clientset[uuid].AppsV1().Deployments(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to collect statistics on the status distribution of the deployment instance,reason: " + err.Error())
		return nil, err
	}
	count := 0
	for _, v := range deployList.Items {
		//readyReplicas 与 replicas 相等 就表示所有pod 都处于ready状态
		if v.Status.ReadyReplicas == *v.Spec.Replicas {
			count++
		}
	}
	return &CountDeploy{
		Total:    len(deployList.Items),
		Ready:    count,
		NotReady: len(deployList.Items) - count,
	}, nil
}

/*

type DeploymentCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	ContainerName string            `json:"container_name"`
	Image         string            `json:"image"`
	Labels        map[string]string `json:"labels"`
	Cpu           string            `json:"cpu" `
	Memory        string            `json:"memory" `
	PortName      string            `json:"port_name"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Protocol      corev1.Protocol   `json:"protocol"`
}
deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:   data.Name,
			Labels: data.Labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Labels,
				},
				Spec: corev1.PodSpec{
					Volumes:        nil,
					InitContainers: nil,
					Containers: []corev1.Container{
						{
							Name:  data.ContainerName,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          data.PortName,
									ContainerPort: data.ContainerPort,
									Protocol:      data.Protocol,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU:    resource.MustParse(data.Cpu),
									corev1.ResourceMemory: resource.MustParse(data.Memory),
								},
								Requests: map[corev1.ResourceName]resource.Quantity{
									corev1.ResourceCPU:    resource.MustParse(data.Cpu),
									corev1.ResourceMemory: resource.MustParse(data.Memory),
								},
							},
						},
					},
				},
			},
		},
	}
	// 判断是否打开健康检测
	if data.HealthCheck {
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			//设置第一个容器的readinessprobe
			// 若存在多个容器 使用for 进行定义
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						// intstr.IntOrString 的作用是端口可以定义为整型，也可以定义为字符
						// type=0 表示该结构体实例内的数据为整型，转json 时只需要使用 IntVal 的数据
						// type=1 表示该结构体实例内的数据为字符串，转json 时只需要使用 StrVal 的数据
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 15,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
*/
