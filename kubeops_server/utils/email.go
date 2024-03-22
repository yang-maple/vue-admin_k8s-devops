package utils

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"kubeops/config"
	"path/filepath"
	"strings"
	"time"
)

// EmailServer 邮件服务端
var EmailServer *mail.SMTPServer

func InitEmail() {
	//创建服务端邮件客户端
	Server := mail.NewSMTPClient()
	//配置客户端信息
	Server.Host = config.HostSmtp
	Server.Port = config.Port
	Server.Username = config.Username
	Server.Password = config.Password
	Server.Encryption = mail.EncryptionNone
	Server.ConnectTimeout = 10 * time.Second
	Server.SendTimeout = 10 * time.Second
	// 客户端赋值
	EmailServer = Server
	client, err := EmailServer.Connect()
	if err != nil {
		Logger.Fatal("The mailbox connected was failed " + err.Error())
	}
	defer client.Close()
	Logger.Info("The mailbox connected was successful")

}

// Emails 发送邮件
func Emails(emails string, content string, subject string) (err error) {
	//创建服务端邮件客户端
	//EmailServer := mail.NewSMTPClient()
	////配置客户端信息
	//EmailServer.Host = config.HostSmtp
	//EmailServer.Port = config.Port
	//EmailServer.Username = config.Username
	//EmailServer.Password = config.Password
	//EmailServer.Encryption = mail.EncryptionNone
	//EmailServer.ConnectTimeout = 10 * time.Second
	//EmailServer.SendTimeout = 10 * time.Second
	//EmailServer.KeepAlive = false
	//EmailServer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	//创建客户端
	EmailClient, err := EmailServer.Connect()
	if err != nil {
		Logger.Fatal("The mailbox connected was failed " + err.Error())
	}
	email := mail.NewMSG()

	email.SetFrom(config.Address).
		AddTo(emails).
		//AddCc("otherto@example.com").
		SetSubject(subject)

	email.SetBody(mail.TextHTML, content)
	err = email.Send(EmailClient)
	if err != nil {
		return err
	}
	//发送成功后关闭连接
	defer EmailClient.Close()
	return nil
}

// FormatEmailBody 格式化邮件内容
func FormatEmailBody(path string, data interface{}) string {
	builder := &strings.Builder{}
	funcs := map[string]interface{}{
		"dateformat": func(t *time.Time) string {
			if t == nil {
				return ""
			}
			return t.Format("2006-01-02 15:04:05")
		},
	}
	tpl := template.Must(template.New("tpl").Funcs(funcs).ParseFiles(path))
	err := tpl.ExecuteTemplate(builder, filepath.Base(path), data)
	if err != nil {
		return err.Error()
	}
	return builder.String()
}

// TestEmail 测试连通性
func TestEmail() {
	//创建服务端邮件客户端
	//EmailServer := mail.NewSMTPClient()
	////配置客户端信息
	//EmailServer.Host = config.HostSmtp
	//EmailServer.Port = config.Port
	//EmailServer.Username = config.Username
	//EmailServer.Password = config.Password
	//EmailServer.Encryption = mail.EncryptionNone
	//EmailServer.ConnectTimeout = 10 * time.Second
	//EmailServer.SendTimeout = 10 * time.Second
}
