package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"go-proxy-auth-service/internal/env"
	"go-proxy-auth-service/internal/jwt"
	"go-proxy-auth-service/internal/middleware"
	"go-proxy-auth-service/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	config := env.GetEnv()
	jwt.GetJwtSecret()
	if config.DisableLogs != 0 {
		log.SetOutput(io.Discard)
		gin.DefaultWriter = io.Discard
	}
	db, err := gorm.Open(sqlite.Open(config.DbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}
	repos := repositories.NewRepository(db)
	ginCtx := gin.Default()
	target, _ := url.Parse("http://127.0.0.1:3000")

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set("X-Proxy", "gin")
	}

	log.Printf("Listening on %v, proxies to %v\n", config.Port, config.TargetUrl)
	ginCtx.Any("/*path", middleware.AuthMiddleware(repos), gin.WrapH(proxy))
	ginCtx.Run(":" + config.Port)
}
