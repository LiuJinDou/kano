package middleware

import (
	"github.com/gin-gonic/gin"
)

func LoginAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authentication := c.Request.Header.Get("Authentication")
		// if len(authentication) == 0 {
		// 	logger.Error("token 不存在")
		// 	tool.SetResponseError(c, http.StatusUnauthorized, errors.New("token 不存在"))
		// 	c.Abort()
		// 	return
		// }

		// token, err := base64.StdEncoding.DecodeString(authentication)
		// if err != nil {
		// 	logger.Error("base64 解析 token 失败, token 无效: ", err)
		// 	tool.SetResponseError(c, http.StatusUnauthorized, err)
		// 	c.Abort()
		// 	return
		// }
		// // 解析 ticket
		// tokenData, err := tool.ParseToken(string(token))
		// if err != nil {
		// 	logger.Error("jwt 解析 token 失败, token 无效: ", err)
		// 	tool.SetResponseError(c, http.StatusUnauthorized, err)
		// 	c.Abort()
		// 	return
		// }

		// logger.Info(tokenData.Username)

		// c.Set("username", tokenData.Username)
		c.Set("username", "dev")
		// c.Set("account_level", tokenData.AccountLevel)

		c.Next()
	}
}

func GetUserName(c *gin.Context) (string, error) {
	if value, exists := c.Get("username"); exists {
		return value.(string), nil
	}
	return "", nil
}
