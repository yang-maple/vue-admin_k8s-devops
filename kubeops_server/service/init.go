package service

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubeops/utils"
)

type k8s struct {
	Clientset map[int]*kubernetes.Clientset
	ConfigDir map[int]*string
}

var K8s = &k8s{
	Clientset: make(map[int]*kubernetes.Clientset),
	ConfigDir: make(map[int]*string),
}

// Init 初始化方法
// func (k *k8s) Init(dir string, uuid int) error {
// 是否初始化 - 是(关闭连接，重新连接) - 否(直接连接)
// 关闭连接存在问题
// 读取配置文件
func (k *k8s) Init(uuid int) error {
	conf, err := clientcmd.BuildConfigFromFlags("", *K8s.ConfigDir[uuid])
	if err != nil {
		utils.Logger.Error("Failed to get configuration file,reason: " + err.Error())
		return err
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		utils.Logger.Error("Failed to create client set,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("k8s initialization succeeded")
	k.Clientset[uuid] = clientset
	return nil
}
