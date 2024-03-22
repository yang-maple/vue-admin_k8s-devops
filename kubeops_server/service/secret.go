package service

import (
	"context"
	"encoding/base64"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
)

type secret struct{}
type SecretResp struct {
	Total int          `json:"total"`
	Item  []secretInfo `json:"item"`
}

type secretInfo struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
	Immutable *bool             `json:"immutable"`
	Type      corev1.SecretType `json:"type"`
	Age       string            `json:"age"`
}

type CreateSecret struct {
	Name           string            `json:"name"`
	Namespace      string            `json:"namespace"`
	Labels         map[string]string `json:"labels"`
	Immutable      bool              `json:"immutable"`
	Type           string            `json:"type"`
	Data           map[string]string `json:"data"`
	TlsCrt         string            `json:"tls_crt"`
	TlsKey         string            `json:"tls_key"`
	Username       string            `json:"username"`
	Password       string            `json:"password"`
	DockerUsername string            `json:"docker_username"`
	DockerPassword string            `json:"docker_password"`
	DockerEmail    string            `json:"docker_email"`
	DockerRegistry string            `json:"docker_registry"`
}

var Secrets secret

func (s *secret) toCells(secrets []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(secrets))
	for i := range secrets {
		cells[i] = secretCell(secrets[i])
	}
	return cells
}

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}
	return secrets
}

// GetSecretList 列表
func (s *secret) GetSecretList(secretName, Namespace string, Limit, Page int, uuid int) (DP *SecretResp, err error) {
	//获取deployment 的所有清单列表
	secretList, err := K8s.Clientset[uuid].CoreV1().Secrets(Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Secrets list,reason: " + err.Error())
		return nil, err
	}

	//组装数据
	selectData := &dataselector{
		GenericDataList: s.toCells(secretList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{secretName},
			Paginate: &PaginateQuery{
				limit: Limit,
				page:  Page,
			},
		},
	}
	//先过滤 后排序
	filtered := selectData.Filter()
	total := len(filtered.GenericDataList)
	//排序并分页
	dataPage := filtered.Sort().Pagination()
	secrets := s.fromCells(dataPage.GenericDataList)
	item := make([]secretInfo, 0, total)
	for _, v := range secrets {
		item = append(item, secretInfo{
			Name:      v.Name,
			Namespace: v.Namespace,
			Labels:    v.Labels,
			Immutable: v.Immutable,
			Type:      v.Type,
			Age:       v.CreationTimestamp.Time.Format("2006-01-02 15:04:05"),
		})
	}
	return &SecretResp{
		Total: total,
		Item:  item,
	}, nil
}

// GetSecretDetail 详情
func (s *secret) GetSecretDetail(Namespace, secretName string, uuid int) (detail *corev1.Secret, err error) {
	//获取deploy
	detail, err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Secrets " + secretName + " detail,reason: " + err.Error())
		return nil, err
	}
	detail.Kind = "Secret"
	detail.APIVersion = "v1"
	utils.Logger.Info("Get Secrets " + secretName + "success")
	return detail, nil
}

// CreateSecret 创建
func (s *secret) CreateSecret(data *CreateSecret, uuid int) (err error) {
	var newsecret *corev1.Secret
	fmt.Println(data.Type)
	switch data.Type {
	case "Opaque":
		newsecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: data.Namespace,
			},
			Type:       corev1.SecretTypeOpaque,
			StringData: data.Data,
		}
	case "kubernetes.io/dockerconfigjson":
		authConfig := fmt.Sprintf(`{"auths":{"%s":{"username":"%s","password":"%s","email":"%s","auth":"%s"}}}`, data.DockerRegistry, data.DockerUsername, data.DockerPassword, data.DockerEmail, base64.StdEncoding.EncodeToString([]byte(data.DockerUsername+":"+data.DockerPassword)))
		newsecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: data.Namespace,
			},
			Type: corev1.SecretTypeDockerConfigJson,
			Data: map[string][]byte{
				//存放字符串不是存放base64编码
				".dockerconfigjson": []byte(authConfig),
			},
		}
	case "kubernetes.io/tls":
		newsecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: data.Namespace,
			},
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSCertKey:       []byte(data.TlsCrt),
				corev1.TLSPrivateKeyKey: []byte(data.TlsKey),
			},
		}
	case "kubernetes.io/basic-auth":
		newsecret = &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: data.Namespace,
			},
			Type: corev1.SecretTypeBasicAuth,
			StringData: map[string]string{
				"username": data.Username,
				"password": data.Password,
			},
		}

	}

	_, err = K8s.Clientset[uuid].CoreV1().Secrets(data.Namespace).Create(context.TODO(), newsecret, metav1.CreateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Create the Secrets " + data.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Create Secrets " + data.Name + "success")
	return nil

}

// DelSecret 删除
func (s *secret) DelSecret(Namespace, secretName string, uuid int) (err error) {
	err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Delete the Secrets " + secretName + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Delete Secrets " + secretName + "success")
	return nil
}

// UpdateSecret   更新
func (s *secret) UpdateSecret(Namespace string, Config *corev1.Secret, uuid int) (err error) {
	_, err = K8s.Clientset[uuid].CoreV1().Secrets(Namespace).Update(context.TODO(), Config, metav1.UpdateOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Update the Secrets " + Config.Name + " ,reason: " + err.Error())
		return err
	}
	utils.Logger.Info("Update Secrets " + Config.Name + "success")
	return nil
}
