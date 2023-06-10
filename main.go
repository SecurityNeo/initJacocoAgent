package main

import (
	"crypto/tls"
	"flag"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"initJacocoAgent/controller"
	"initJacocoAgent/middleware"
	"initJacocoAgent/types"
	"strconv"
	"strings"
	"time"
)

var Param types.WHSrvParam

func init() {
	flag.IntVar(&Param.Port, "port", 443, "webhook server port")
	flag.StringVar(&Param.CertFile, "tlsCertFile", "/opt/certs/cert.pem", "certificate file for https")
	flag.StringVar(&Param.KeyFile, "tlsKeyFile", "/opt/certs/key.pem", "private key file")
	flag.Var(&Param.ProtectNS, "protect_ns", "protect namespaces")
	flag.Parse()
}

func handler(handler types.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := new(types.NewContext)
		ctx.Context = c
		ctx.ProtectNS = Param.ProtectNS
		handler(ctx)
	}
}

func main() {
	_, err := tls.LoadX509KeyPair(Param.CertFile, Param.KeyFile)
	if err != nil {
		log.Errorf("Failed to load key pair: %v", err)
	}
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(middleware.Logger(), gin.Recovery())

	v1Group := router.Group("api/v1")
	{
		v1Group.POST("/mutate", handler(controller.WebhookCallBack))
	}
	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 30 * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20

	srvAddr := ":" + strconv.Itoa(Param.Port)
	srv := endless.NewServer(srvAddr, router)

	err = srv.ListenAndServeTLS(Param.CertFile, Param.KeyFile)
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Error("Listen: %s\n", err)
		}
	}
}
