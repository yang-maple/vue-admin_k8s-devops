package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	yamlv3 "gopkg.in/yaml.v3"
	"io"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"kubeops/utils"
	"os"
)

var Upload upload

type upload struct{}

// 定义资源结构体，用于返回创建的资源对象
type fileResource struct {
	kind      string
	name      string
	namespace string
}

// UploadFile 通过上传yaml文件创建资源
func (u *upload) UploadFile(dst string, uuid int) (msg string, err error) {
	//读取文件内容
	filebytes, err := os.ReadFile(dst)
	if err != nil {
		utils.Logger.Error("Failed to read the file,reason: " + err.Error())
		return "", err
	}
	//传入文件内容，创建资源
	msg, err = Upload.CreateYaml(filebytes, uuid)
	if err != nil {
		utils.Logger.Error("Failed to create a resource,reason: " + err.Error())
		return "", err
	}
	//执行完成删除文件
	_ = os.Remove(dst)
	// 获取文件后缀
	return msg, nil
}

// CreateYaml 通过传入yaml内容创建资源
func (u *upload) CreateYaml(fileBytes []byte, uuid int) (msg string, err error) {
	// yaml 格式校验
	var result map[string]interface{}
	err = yamlv3.Unmarshal(fileBytes, &result)
	if err != nil {
		utils.Logger.Error("The yaml format is invalid" + err.Error())
		return "", errors.New("yaml 格式无效")
	}
	// 定义资源变量
	var resource fileResource
	//默认值
	nameSpace := "default"
	// 创建客户端连接
	conf, err := clientcmd.BuildConfigFromFlags("", *K8s.ConfigDir[uuid])
	if err != nil {
		utils.Logger.Error("Failed to create a ConfigClient" + err.Error())
		return "", err
	}
	dd, err := dynamic.NewForConfig(conf)
	if err != nil {
		utils.Logger.Error("Failed to create a DynamicClient" + err.Error())
		return "", err
	}
	//缓存100 读取数据
	decoder := yamlutil.NewYAMLOrJSONDecoder(bytes.NewReader(fileBytes), 100)
	for {
		//读取yaml文件 读取完成时返回资源对象
		var rawObj runtime.RawExtension
		if err = decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				return fmt.Sprintf("%s/%s created in %s", resource.kind, resource.name, resource.namespace), nil
			}
			return "", err
		}
		// 资源转化
		obj, gvk, err := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode(rawObj.Raw, nil, nil)
		unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
		if err != nil {
			utils.Logger.Error("The conversion failed" + err.Error())
			return "", errors.New("转换失败" + err.Error())
		}
		// 获取APIGroupResources 用于获取RESTMapping
		unstructuredObj := &unstructured.Unstructured{Object: unstructuredMap}
		gr, err := restmapper.GetAPIGroupResources(K8s.Clientset[uuid].Discovery())
		if err != nil {
			utils.Logger.Error("Failed to get APIGroupResources " + err.Error())
			return "", errors.New("获取APIGroupResources失败" + err.Error())
		}
		// 获取RESTMapping 用于获取ResourceInterface
		mapper := restmapper.NewDiscoveryRESTMapper(gr)
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			utils.Logger.Error("Failed to get RESTMapping" + err.Error())
			return "", errors.New("获取RESTMapping失败" + err.Error())
		}
		// 使用 ResourceInterface 创建资源
		var dri dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			if unstructuredObj.GetNamespace() == "" {
				unstructuredObj.SetNamespace(nameSpace)
			}
			dri = dd.Resource(mapping.Resource).Namespace(unstructuredObj.GetNamespace())
		} else {
			dri = dd.Resource(mapping.Resource)
		}
		obj2, err := dri.Create(context.Background(), unstructuredObj, metav1.CreateOptions{})
		if err != nil {
			utils.Logger.Error("Failed to create a resource" + err.Error())
			return "", err
		}
		resource.kind = obj2.GetKind()
		resource.name = obj2.GetName()
		resource.namespace = obj2.GetNamespace()
		utils.Logger.Info(fmt.Sprintf(" %s / %s created in %s", obj2.GetKind(), obj2.GetName(), obj2.GetNamespace()))
	}
}
