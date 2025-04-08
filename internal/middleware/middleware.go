package middleware

import (
	"github.com/gin-gonic/gin"
	"go-proxy-auth-service/internal/env"
	"go-proxy-auth-service/internal/jwt"
	"go-proxy-auth-service/internal/repositories"
	"net/http"
)

func AuthMiddleware(repos *repositories.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		config := env.GetEnv()
		accessToken, err := c.Cookie(config.AccessTokenName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error(), "code": http.StatusUnauthorized})
			return
		}

		tokenResult, errResult := jwt.VerifyToken(accessToken)
		if errResult != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": errResult.Error(), "code": http.StatusUnauthorized})
		}

		user, err := repos.FindUserWithSession(tokenResult.Login, tokenResult.SessionId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "user/session not found", "code": http.StatusUnauthorized})
		}
		if user.Locked != 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login is locked", "code": http.StatusUnauthorized})
		}
		if user.Role == "operator" && config.RoleOperatorBlock == 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "operator role is prohibited", "code": http.StatusUnauthorized})
		}
		if user.Role == "guest" && config.RoleGuestBlock == 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "guest role is prohibited", "code": http.StatusUnauthorized})
		}
		if user.Role == "admin" && config.RoleAdminBlock == 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "admin role is prohibited", "code": http.StatusUnauthorized})
		}
	}
}
