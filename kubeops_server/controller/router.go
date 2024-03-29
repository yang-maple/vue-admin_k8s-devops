package controller

import (
	"github.com/gin-gonic/gin"
)

// 定义 router 结构体
type router struct{}

// Router 实例化结构体
var Router router

// InitApiRouter 初始化路由规则
func (router *router) InitApiRouter(r *gin.Engine) {
	podV1 := r.Group("/v1/api/pod")
	{
		podV1.GET("/list", Pod.GetPods)
		podV1.DELETE("/delete", Pod.DeletePod)
		podV1.PUT("/update", Pod.UpdatePod)
		podV1.GET("/container/list", Pod.GetContainerList)
		podV1.GET("/detail", Pod.GetPodsDetail)
		podV1.GET("/container/log", Pod.GetContainerLog)
		podV1.GET("/ns", Pod.GetPodNumber)
	}
	deployV1 := r.Group("/v1/api/deployment")
	{
		deployV1.GET("/list", Deployment.GetDeploylist)
		deployV1.PUT("/modify", Deployment.ModifyDeployReplicas)
		deployV1.POST("/restart", Deployment.RestartDeploy)
		deployV1.GET("/detail", Deployment.GetDeployDetail)
		deployV1.POST("/create", Deployment.CreateDeploy)
		deployV1.DELETE("/delete", Deployment.DelDeploy)
		deployV1.PUT("/update", Deployment.UpdateDeploy)
		deployV1.GET("/ns", Deployment.GetDeployPer)
		deployV1.POST("/rollout", Deployment.RolloutDeploy)
	}
	daemonV1 := r.Group("/v1/api/daemon")
	{
		daemonV1.GET("/list", DaemonSet.GetDaemonList)
		daemonV1.GET("/detail", DaemonSet.GetDaemonDetail)
		daemonV1.DELETE("/delete", DaemonSet.DelDaemon)
		daemonV1.PUT("/update", DaemonSet.UpdateDaemon)
	}
	statefulV1 := r.Group("/v1/api/stateful")
	{
		statefulV1.GET("/list", StatefulSet.GetStatefulSetList)
		statefulV1.GET("/detail", StatefulSet.GetStatefulDetail)
		statefulV1.DELETE("/delete", StatefulSet.DelStatefulSet)
		statefulV1.PUT("/update", StatefulSet.UpDataStatefulSet)
		statefulV1.PUT("/modify", StatefulSet.ModifyStatefulReplicas)
	}
	namespacesV1 := r.Group("/v1/api/namespaces")
	{
		namespacesV1.GET("/list", Namespace.GetNsList)
		namespacesV1.GET("/detail", Namespace.GetNsDetail)
		namespacesV1.DELETE("/delete", Namespace.DelNs)
		namespacesV1.POST("/create", Namespace.CreateNs)
		namespacesV1.PUT("/update", Namespace.UpdateNs)
	}
	nodeV1 := r.Group("/v1/api/node")
	{
		nodeV1.GET("/list", Node.GetNodeList)
		nodeV1.GET("/detail", Node.GetNodeDetail)
		nodeV1.POST("/schedule", Node.SetNodeSchedule)
	}
	pvV1 := r.Group("/v1/api/pv")
	{
		pvV1.GET("/list", Persistentvolume.GetPersistentVolumeList)
		pvV1.GET("/detail", Persistentvolume.GetPersistentVolumeDetail)
		pvV1.DELETE("/delete", Persistentvolume.DelPersistentVolume)
		pvV1.POST("/create", Persistentvolume.CreatePersistentVolume)
		pvV1.PUT("/update", Persistentvolume.UpdatePersistentVolume)
	}
	pvcV1 := r.Group("/v1/api/pvc")
	{
		pvcV1.GET("/list", PersistentVolumeClaim.GetPersistentVolumeClaimList)
		pvcV1.GET("/detail", PersistentVolumeClaim.GetPersistentVolumeClaimDetail)
		pvcV1.DELETE("/delete", PersistentVolumeClaim.DelPersistentVolumeClaim)
		pvcV1.POST("/create", PersistentVolumeClaim.CreatePersistentVolumeClaim)
		pvcV1.PUT("/update", PersistentVolumeClaim.UpdatePersistentVolumeClaim)
	}
	svcV1 := r.Group("/v1/api/svc")
	{
		svcV1.GET("/list", Services.GetServiceList)
		svcV1.GET("/detail", Services.GetServiceDetail)
		svcV1.DELETE("/delete", Services.DelServices)
		svcV1.POST("/create", Services.CreateService)
		svcV1.PUT("/update", Services.UpdateService)
	}
	ingV1 := r.Group("/v1/api/ing")
	{
		ingV1.GET("/list", Ingress.GetIngressList)
		ingV1.GET("/detail", Ingress.GetIngressDetail)
		ingV1.DELETE("/delete", Ingress.DelIngress)
		ingV1.POST("/create", Ingress.CreateIngress)
		ingV1.PUT("/update", Ingress.UpdateIngress)
	}
	cmV1 := r.Group("/v1/api/cm")
	{
		cmV1.GET("/list", Configmaps.GetConfigmapList)
		cmV1.GET("/detail", Configmaps.GetConfigmapDetail)
		cmV1.DELETE("/delete", Configmaps.DelConfigmap)
		cmV1.POST("/create", Configmaps.CreateConfigmap)
		cmV1.PUT("/update", Configmaps.UpdateConfigmap)
	}
	secretV1 := r.Group("/v1/api/secret")
	{
		secretV1.GET("/list", Secrets.GetSecretList)
		secretV1.GET("/detail", Secrets.GetSecretDetail)
		secretV1.DELETE("/delete", Secrets.DelSecret)
		secretV1.POST("/create", Secrets.CreateSecret)
		secretV1.PUT("/update", Secrets.UpdateSecret)
	}
	workflowV1 := r.Group("/v1/api/workflow")
	{
		workflowV1.GET("/list", Workflow.GetWorkflowList)
		workflowV1.GET("/detail", Workflow.GetWorkflowDetail)
		workflowV1.DELETE("/delete", Workflow.DeleteWorkflow)
		workflowV1.POST("/create", Workflow.CreateWorkflow)
	}
	userV1 := r.Group("/v1/api/user")
	{
		userV1.GET("/getCaptcha", Login.CaptchaImage)
		userV1.GET("/info", Login.getUserInfo)
		userV1.POST("/login", Login.VerifyInfo)
		userV1.POST("/register", Register.RegisterUser)
		userV1.POST("/register/email", Register.SendEmail)
		userV1.POST("/findPassword/email", ResetPassword.verifyIdentity)
		userV1.POST("/findPassword", ResetPassword.FindPass)
		userV1.POST("/resetPassword", ResetPassword.ResetPass)
	}
	uploadV1 := r.Group("/v1/api/upload")
	{
		uploadV1.POST("/uploadFile", Upload.uploadYamlFile)
		uploadV1.POST("/uploadYaml", Upload.createYaml)
	}
	clusterV1 := r.Group("/v1/api/cluster")
	{
		clusterV1.POST("/create", Cluster.Create)
		clusterV1.GET("/list", Cluster.List)
		clusterV1.POST("/change", Cluster.Change)
		clusterV1.DELETE("/delete", Cluster.Delete)
		clusterV1.PUT("/update", Cluster.Update)
		clusterV1.GET("/detail", Cluster.Get)
	}
	homepageV1 := r.Group("/v1/api/homepage")
	{
		homepageV1.GET("/getInfo", HomePage.GetHomepage)
	}

	ingressClassV1 := r.Group("/v1/api/ingressClass")
	{
		ingressClassV1.GET("/list", IngressClass.GetIngressClassList)
		ingressClassV1.GET("/detail", IngressClass.GetIngressClassDetail)
		ingressClassV1.DELETE("/delete", IngressClass.DelIngressClass)
		ingressClassV1.PUT("/update", IngressClass.UpdateIngressClass)
	}
	storageClassV1 := r.Group("/v1/api/storageClass")
	{
		storageClassV1.GET("/list", StorageClass.GetStorageClassList)
		storageClassV1.GET("/detail", StorageClass.GetStorageClassDetail)
		storageClassV1.DELETE("/delete", StorageClass.DelStorageClass)
		storageClassV1.PUT("/update", StorageClass.UpdateStorageClass)
	}
}
