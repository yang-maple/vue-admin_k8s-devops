package middle

import (
	"github.com/gin-gonic/gin"
	"kubeops/utils"
	"net/http"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//对登陆接口进行放行
		if (len(c.Request.URL.String()) >= 10) && (c.Request.URL.String()[0:12] == "/v1/api/user") {
			c.Next()
		} else {
			token := c.Request.Header.Get("Authorization")
			if token == "" {
				c.JSON(http.StatusOK, gin.H{
					"code": 4001,
					"msg":  "请求未携带token 无权限访问",
					"data": nil,
				})
				c.Abort()
				return
			}
			//解析token 内容
			claims, err := utils.JWTToken.ParseToken(token, utils.UserSecret)
			if err != nil {
				//token 延期错误
				if err.Error() == "TokenExpired" {
					c.JSON(http.StatusOK, gin.H{
						"code": 4002,
						"msg":  "Token 已过期",
						"data": nil,
					})
					c.Abort()
					return
				}
				// 其他解析错误
				c.JSON(http.StatusOK, gin.H{
					"code": 4003,
					"msg":  err.Error(),
					"data": nil,
				})
				c.Abort()
				return
			}
			//继续交由下一个路由进行处理，并将解析出的信息进行传递
			c.Set("claims_id", claims.Id)
			c.Next()
		}

	}
}
