package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"roomino/ctl"
	"roomino/e"
	"roomino/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		token := c.GetHeader("Authorization")
		if token == "" {
			code = http.StatusNotFound
			c.JSON(e.InvalidParams, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Missing Token",
			})
			c.Abort()
			return
		}

		claims, err := util.ParseToken(token)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		}

		if code != e.SUCCESS {
			c.JSON(e.InvalidParams, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Token might have expired. Please log in again.",
			})
			c.Abort()
			return
		}

		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{UserName: claims.Username}))
		c.Next()
	}
}
