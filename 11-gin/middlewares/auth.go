package middlewares

import `github.com/gin-gonic/gin`

func MyAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"hliangzhao": "passwd",
	})
}
