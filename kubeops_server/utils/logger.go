package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"kubeops/config"
	"path"
	"runtime"
	"strings"
	"time"
)

var Logger *logrus.Logger

type MyFormatter struct {
}

// Format 实现 Formatter 接口 Format方法返回自定义日志格式
func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	caller := getCallerInfo()
	if strings.Contains(caller.File, "github.com/sirupsen/logrus@v1.9.3/logger.go") {
		return []byte(fmt.Sprintf("%s [%s] [GIN] %s\n", entry.Time.Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), entry.Message)), nil
	}
	return []byte(fmt.Sprintf("%s [%s] [%s:%d] %s\n", entry.Time.Format("2006-01-02 15:04:05"), strings.ToUpper(entry.Level.String()), caller.File, caller.Line, entry.Message)), nil
}

// LoggerToFile 定义日志中间件格式以及输出
func LoggerToFile() gin.HandlerFunc {
	logFilePath := config.Log_FILE_PATH
	logFileName := config.LOG_FILE_NAME
	//日志文件
	fileName := path.Join(logFilePath, logFileName)

	//写入文件
	//src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	//if err != nil {
	//	fmt.Println("err", err)
	//}

	//实例化
	Logger = logrus.New()

	//日志分割
	writer, _ := rotatelogs.New(
		logFilePath+"%Y-%m-%d-"+logFileName,
		rotatelogs.WithLinkName(fileName),
		rotatelogs.WithRotationCount(10),
		rotatelogs.WithRotationTime(time.Minute*1),
	)
	//设置输出
	Logger.SetOutput(writer)
	//Logger.Out = src

	//设置日志级别
	Logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	Logger.SetFormatter(new(MyFormatter))

	//设置日志分割

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 请求用户
		clientUserId := c.Request.Header.Get("Uuid")

		// 日志格式
		Logger.Infof("UserId = %s | %3d | %13v | %15s | %s | %s |",
			clientUserId,
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// getCallerInfo 获取调用者信息
func getCallerInfo() *runtime.Frame {
	pc, file, line, ok := runtime.Caller(7) // 通过runtime.Caller获取调用者信息
	if !ok {
		return nil
	}
	// 获取调用者信息
	caller := runtime.FuncForPC(pc)
	// 获取文件名和行号
	return &runtime.Frame{Function: path.Base(caller.Name()), File: file, Line: line}
}
