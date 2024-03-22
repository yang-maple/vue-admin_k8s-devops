package controller

import (
	"github.com/ghodss/yaml"
	"github.com/gin-gonic/gin"
	"kubeops/service"
	"kubeops/utils"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

type YamlData struct {
}

var Upload upload

type upload struct{}

const yamlurl = "./static/yaml/"

// 上传yaml文件创建资源
func (u *upload) uploadYamlFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["yamlfile"]
	for _, file := range files {
		// 上传文件至指定目录
		extname := path.Ext(file.Filename)
		exifile := map[string]bool{
			".yaml": true,
		}
		if exifile[extname] {
			dir := yamlurl + time.Now().Format("2006-01-02")
			if !utils.DirExists(dir) {
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"code": 400,
						"msg":  "文件上传失败",
					})
					return
				}
			}
			dst := path.Join(dir, file.Filename)
			err := c.SaveUploadedFile(file, dst)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "文件上传失败" + err.Error(),
				})
				return
			}
			uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
			msg, err := service.Upload.UploadFile(dst, uuid)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"code": 400,
					"msg":  "资源创建失败" + err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 200,
				"msg":  msg,
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "上传文件格式不正确，请上传yaml文件",
		})
	}
}

// 通过 yaml 内容创建资源
func (u *upload) createYaml(c *gin.Context) {
	//定义并接受yaml 内容
	params := new(struct {
		YamlContent string `json:"yamlContent"`
	})
	_ = c.ShouldBindJSON(&params)
	// 将json 格式转化为yaml
	yamlStr, err := yaml.JSONToYAML([]byte(params.YamlContent))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "yaml 格式无效",
		})
		return
	}
	// 创建资源
	uuid, _ := strconv.Atoi(c.Request.Header.Get("Uuid"))
	msg, err := service.Upload.CreateYaml(yamlStr, uuid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "资源创建失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
	})
}
