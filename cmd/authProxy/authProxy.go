// @title           auth proxy API
// @version         1.0
// @description     This is API with Swagger documentation.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  monkeyhouse@mail.ru

// @host      localhost:9080
// @BasePath  /

package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	_ "go-proxy-auth-service/docs"
	"go-proxy-auth-service/internal/env"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := env.GetEnv()
		accessToken, err := c.Cookie(config.AccessTokenName)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		log.Printf("req %s %s\n", c.Request.URL.String(), accessToken)
		//c.AbortWithStatus(http.StatusForbidden)
	}
}
func main() {
	config := env.GetEnv()
	if config.DisableLogs != 0 {
		log.SetOutput(io.Discard)
		gin.DefaultWriter = io.Discard
	}
	ginCtx := gin.Default()
	target, _ := url.Parse("http://127.0.0.1:3000")
	//target, _ := url.Parse("http://localhost:9000")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set("X-Proxy", "gin")
	}
	// Middleware для дополнительных проверок
	//ginCtx.Use(func(c *gin.Context) {
	// Проверка перед проксированием
	/*if c.Request.URL.Path == "/admin" && !isAdmin(c) {
		c.AbortWithStatus(403)
		return
	}*/
	//	c.Next()
	//})
	//swagger.Controller(ginCtx)

	log.Printf("Listening on %v, proxies to %v\n", config.Port, config.TargetUrl)
	ginCtx.Any("/*path", AuthMiddleware(), gin.WrapH(proxy))
	ginCtx.Run(":" + config.Port)
}
