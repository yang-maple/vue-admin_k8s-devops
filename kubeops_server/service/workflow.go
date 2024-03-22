package service

import (
	"kubeops/dao"
	"kubeops/model"
)

//workflow 用于处理工作流

type workflow struct {
}

var Workflow workflow

//	type Workflowcreate struct {
//		Name          string            `json:"name"`
//		Namespace     string            `json:"namespace"`
//		Replicas      int32             `json:"replicas"`
//		Image         string            `json:"image"`
//		Labels        map[string]string `json:"labels"`
//		Cpu           string            `json:"cpu"`
//		Memory        string            `json:"memory"`
//		ContainerPort int32             `json:"container_port"`
//		HealthCheck   bool              `json:"health_check"`
//		HealthPath    string            `json:"health_path"`
//		Type          string            `json:"type"`
//		Port          int32             `json:"port"`
//		NodePort      int32             `json:"node_port"`
//		Hosts         []Rule            `json:"hosts"`
//	}
type Workflowcreate struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Type      string            `json:"type"`
	Deploy    DeploymentCreate  `json:"deploy"`
	Service   CreateService     `json:"service"`
	Ingress   CreateIngress     `json:"ingress"`
}

// GetWorkflowList 获取工作流列表
func (w *workflow) GetWorkflowList(FilterName, Namespace string, Limit, Page int) (workflow *dao.WorkflowResp, err error) {
	data, err := dao.Workflow.GetWorkflowList(FilterName, Namespace, Limit, Page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetWorkflowDetail 获取工作流详细
func (w *workflow) GetWorkflowDetail(id int) (workflow *model.Workflow, err error) {
	data, err := dao.Workflow.GetByIdDetail(id)
	if err != nil {
		return nil, err
	}

	return data, err
}

// CreateWorkflow 新增工作流 workflow
func (w *workflow) CreateWorkflow(data *Workflowcreate, uuid int) (err error) {
	//定义一个 ingressname 用于判断 该实例是否 用到 ingress 类型
	var ingressname string
	//判断 该实例类型是否为 ingress
	if data.Type == "Ingress" {
		ingressname = Getingressname(data.Name)
	} else {
		ingressname = ""
	}

	//组装数据
	workflow := &model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Deploy.Replicas,
		Deployment: data.Name,
		Service:    Getservicename(data.Name),
		Ingress:    ingressname,
		Type:       data.Type,
	}
	//数据写入数据库

	//创建k8s 资源
	err = createWorkflowRes(data, uuid)
	if err != nil {
		return err
	} else {
		//后端创建成功后，数据写入数据库，否则不进行写入操作
		err = dao.Workflow.Add(workflow)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteById 删除工作流
func (w *workflow) DeleteById(id int, uuid int) (err error) {
	//获取数据库资源
	workflow, err := dao.Workflow.GetById(id)
	if err != nil {
		return err
	}

	//删除k8s资源
	err = deleteWorkflowRes(workflow, uuid)
	if err != nil {
		return err
	}
	//删除数据库数据
	err = dao.Workflow.DeleteById(id)
	if err != nil {
		return err
	}
	return nil
}

// 创建k8s 资源的函数 deploy service ingress
func createWorkflowRes(data *Workflowcreate, uuid int) (err error) {
	//组装数据
	deploy := &data.Deploy
	//给公共字段赋值
	deploy.Name = data.Name
	deploy.Namespace = data.Namespace
	deploy.Labels = data.Labels
	//创建deploy
	err = Deployment.CreateDeploy(deploy, uuid)
	if err != nil {
		return err
	}

	//创建 svc 资源
	//判断 service 的类型
	var svctype string
	if data.Type != "Ingress" {
		svctype = data.Type
	} else {
		svctype = "ClusterIP"
	}
	//组装 service 资源
	svc := &data.Service
	//公共字段赋值
	svc.Name = Getservicename(data.Name)
	svc.Namespace = data.Namespace
	svc.Labels = data.Labels
	//类型更换
	svc.Type = svctype
	//创建资源
	err = Services.CreateSvc(svc, uuid)
	if err != nil {
		_ = Deployment.DelDeploy(data.Namespace, data.Name, uuid)
		return err
	}
	//创建资源
	//判断 ingress 的类型
	var ing *CreateIngress
	if data.Type == "Ingress" {
		//组装数据
		ing = &data.Ingress
		//公共字段赋值
		ing.Name = Getingressname(data.Name)
		ing.Namespace = data.Namespace
		ing.Labels = data.Labels
		for i := range ing.Rules {
			for j := range ing.Rules[i].HTTPIngressRuleValues {
				ing.Rules[i].HTTPIngressRuleValues[j].ServiceName = Getservicename(data.Name)
			}
		}
	} else {
		return nil
	}
	if err := Ingress.CreateIng(ing, uuid); err != nil {
		_ = Deployment.DelDeploy(data.Namespace, data.Name, uuid)
		_ = Services.DelSvc(data.Namespace, Getservicename(data.Name), uuid)
		return err
	}
	return nil
}

// 删除 k8s资源的函数 deploy service ingress
func deleteWorkflowRes(data *model.Workflow, uuid int) (err error) {
	//判断是否为ingress 如果是则删除 ingress 资源 如果没有ingress资源则跳过
	if data.Type == "Ingress" {
		err = Ingress.DelIng(data.Namespace, Getingressname(data.Name), uuid)
		if err != nil {
			return err
		}
	}
	//删除 service 资源
	err = Services.DelSvc(data.Namespace, Getservicename(data.Name), uuid)
	if err != nil {
		return err
	}

	//删除 deploy 资源
	err = Deployment.DelDeploy(data.Namespace, data.Name, uuid)
	if err != nil {
		return err
	}
	return nil
}

// Getservicename 获取svc name
func Getservicename(data string) string {
	return data + "svc"
}

// Getingressname  获取ingress name
func Getingressname(data string) string {
	return data + "ing"
}
