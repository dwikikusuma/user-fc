package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleWare(secrete string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"err":     "authorization cannot ben nil",
			})

			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"err":     "invalid authorization format",
			})

			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(secrete), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"err":     "token might be expired",
				"exc":     err,
				"secre":   secrete,
			})

			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"err":     "failed to unpack token",
			})

			c.Abort()
			return
		}

		c.Set("user_id", claims["user_id"].(float64))
		c.Next()
	}
}
