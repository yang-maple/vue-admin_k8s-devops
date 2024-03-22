package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"kubeops/utils"
	"net/http"
	"time"
)

type terminal struct{}

var Terminal terminal

func (t *terminal) WsHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	if err := r.ParseForm(); err != nil {
		utils.Logger.Error("parse failed" + err.Error())
		return
	}
	namespace := r.Form.Get("namespace")
	podName := r.Form.Get("pod_name")
	containerName := r.Form.Get("container_name")
	token := r.Form.Get("token")
	//验证token
	clamis, err := utils.JWTToken.ParseToken(token, utils.UserSecret)
	if err != nil {
		utils.Logger.Error("Token verification failed" + err.Error())
		_, _ = w.Write([]byte("token验证失败" + err.Error()))
		return
	}
	//加载k8s配置
	conf, err := clientcmd.BuildConfigFromFlags("", *K8s.ConfigDir[clamis.Id])
	if err != nil {
		utils.Logger.Error("Failed to load k8s configuration" + err.Error())
		return
	}
	utils.Logger.Info("kubectl exec pod %s -c %s -n %s\n", podName, containerName, namespace)
	pty, err := t.NewTerminalSession(w, r, nil)
	if err != nil {
		utils.Logger.Error("new Terminal session failed" + err.Error())
		return
	}
	defer func() {
		utils.Logger.Info("defer close TerminalSession")
		pty.Close()
	}()
	req := K8s.Clientset[clamis.Id].CoreV1().RESTClient().Post().Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
			Container: containerName,
			Command:   []string{"/bin/bash"},
		}, scheme.ParameterCodec)
	utils.Logger.Info("exec post request url:", req)

	//update SPDY
	executor, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		utils.Logger.Error("connect SPDY failed" + err.Error())
		return
	}

	//
	err = executor.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		Tty:               true,
		TerminalSizeQueue: pty,
	})
	if err != nil {
		utils.Logger.Info("exec cmd failed" + err.Error())
		//return err to web
		_, _ = pty.Write([]byte(("exec cmd failed") + err.Error()))
		// close
		pty.Done()
	}
}

// 定义消息内容的结构体
// Operation 是操作类型
// Data 是传入数据内容
// Rows和Cols 理解为终端的行数和列数
type terminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows" `
	Cols      uint16 `json:"cols" `
}

// TerminalSession 定义 websocket 交互结构体 接管输入和输出
// wsConn 是websocket 连接
// sizeChan 定义 终端输入和输出的宽和高
// doneChan 用于标记退出终端
type TerminalSession struct {
	wsConn   *websocket.Conn
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

// 初始化 websocket.Upgrade 类型的对象,用于http协议升级为 ws 协议
var upgrader = func() websocket.Upgrader {
	upgrader := websocket.Upgrader{}
	upgrader.HandshakeTimeout = time.Second * 2
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return upgrader
}()

// NewTerminalSession 创建terminal 类型的对象并返回
func (t *terminal) NewTerminalSession(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*TerminalSession, error) {
	//升级ws 协议
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, errors.New("update websocket failed" + err.Error())
	}
	//new
	session := &TerminalSession{
		wsConn:   conn,
		sizeChan: make(chan remotecommand.TerminalSize),
		doneChan: make(chan struct{}),
	}
	return session, nil
}

// read data function
func (t *TerminalSession) Read(p []byte) (int, error) {
	// 从conn 中获取message 信息
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		utils.Logger.Error("read message failed" + err.Error())
		return 0, errors.New("read message failed" + err.Error())
	}
	//反序列化 转化成字符串
	var msg terminalMessage
	err = json.Unmarshal(message, &msg)
	if err != nil {
		utils.Logger.Error("json Unmarshal message failed" + err.Error())
		return 0, errors.New("json Unmarshal message failed" + err.Error())
	}
	//逻辑判断
	switch msg.Operation {
	case "stdin":
		msg.Data = msg.Data + "\r"
		return copy(p, msg.Data), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{
			Width:  msg.Cols,
			Height: msg.Rows,
		}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		utils.Logger.Info("unknown message Operation" + msg.Operation)
		return 0, errors.New("unknown message Operation" + err.Error())
	}
}

// writer data function,get api-server return message,send msg to web
func (t *TerminalSession) Write(p []byte) (int, error) {
	msg, err := json.Marshal(terminalMessage{
		Operation: "stdout",
		Data:      string(p),
		Rows:      0,
		Cols:      0,
	},
	)
	if err != nil {
		utils.Logger.Error("write parse message error" + err.Error())
		return 0, errors.New("write parse message error" + err.Error())
	}
	if err := t.wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		utils.Logger.Error("write parse message error" + err.Error())
		return 0, errors.New("write message error" + err.Error())
	}
	return len(p), nil
}

// Done 标记关闭的方法
func (t *TerminalSession) Done() {
	close(t.doneChan)
}

// Close 关闭的方法
func (t *TerminalSession) Close() {
	t.wsConn.Close()
}

// Next resize方法，以及是否退出终端
func (t *TerminalSession) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}
