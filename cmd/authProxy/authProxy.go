package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/golang-jwt/jwt/v5"
	"go-proxy-auth-service/internal/env"
	"go-proxy-auth-service/internal/repositories"
	"go-proxy-auth-service/internal/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type CustomClaims struct {
	Login     string `json:"login"`
	Role      string `json:"role"`
	SessionId string `json:"sessionId"`
	Type      string `json:"type"`
	jwt.RegisteredClaims
}

var jwtSecret *rsa.PublicKey

func parseRSAPublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}

	rsaPub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS1 encoded public key: %w", err)
	}

	return rsaPub, nil
}

func getJwtSecret() {
	config := env.GetEnv()
	data, err := utils.ReadFile(config.PublicKeyPath)
	if err != nil {
		panic("Could not get public key")
	}
	secret, err := parseRSAPublicKey(data)
	if err != nil {
		panic("Could not parse RSA AES public key")
	}
	jwtSecret = secret
}
func VerifyToken(tokenString string) (*CustomClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
func main() {
	config := env.GetEnv()
	getJwtSecret()
	if config.DisableLogs != 0 {
		log.SetOutput(io.Discard)
		gin.DefaultWriter = io.Discard
	}
	db, err := gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	repos := repositories.NewRepository(db)
	fmt.Println("repos", repos)
	ginCtx := gin.Default()
	target, _ := url.Parse("http://127.0.0.1:3000")

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set("X-Proxy", "gin")
	}

	AuthMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			config := env.GetEnv()
			accessToken, err := c.Cookie(config.AccessTokenName)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": err.Error(), "code": http.StatusUnauthorized})
				return
			}

			tokenResult, errResult := VerifyToken(accessToken)
			if errResult != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": errResult.Error(), "code": http.StatusUnauthorized})
			}

			user, err := repos.FindUser(tokenResult.Login, tokenResult.SessionId)

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
	log.Printf("Listening on %v, proxies to %v\n", config.Port, config.TargetUrl)
	ginCtx.Any("/*path", AuthMiddleware(), gin.WrapH(proxy))
	ginCtx.Run(":" + config.Port)
}
