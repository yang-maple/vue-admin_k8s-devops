package service

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeops/utils"
	"sort"
)

var Event event

type event struct {
}

type EventResp struct {
	Total int         `json:"total"`
	Item  []eventInfo `json:"item"`
}

type eventInfo struct {
	Namespace string `json:"namespace"`
	LastTime  string `json:"last_time"`
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Object    string `json:"object"`
	Message   string `json:"message"`
}

func (e *event) ListEvent(uuid int, types string, namespaces string) (info *EventResp, err error) {
	events, err := K8s.Clientset[uuid].CoreV1().Events(namespaces).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		utils.Logger.Error("Failed to Get the Events list,reason: " + err.Error())
		return nil, err
	}
	//定义新的[]info 遍历event 循环写入info
	item := make([]eventInfo, 0, len(events.Items))
	for _, event := range events.Items {
		if event.Type != types {
			item = append(item, eventInfo{
				Namespace: event.Namespace,
				LastTime:  event.LastTimestamp.Format("2006-01-02 15:04:05"),
				Type:      event.Type,
				Reason:    event.Reason,
				Object:    event.InvolvedObject.Name,
				Message:   event.Message,
			})
		}
	}
	sort.Slice(item, func(i, j int) bool {
		return item[i].LastTime > item[j].LastTime
	})
	return &EventResp{
		Total: len(events.Items),
		Item:  item,
	}, nil
}
