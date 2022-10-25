package main

import (
	"crypto/tls"
	"flag"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"initSkywalkingAgent/handler"
	"strconv"
	"strings"
	"time"
)

type WHSrvParam struct {
	port     int
	certFile string
	keyFile  string
}

func main() {

	var param WHSrvParam

	flag.IntVar(&param.port, "port", 1443, "webhook server port")
	flag.StringVar(&param.certFile, "tlsCertFile", "/opt/certs/cert.pem", "certificate file for https")
	flag.StringVar(&param.keyFile, "tlsKeyFile", "/opt/certs/key.pem", "private key file")
	flag.Parse()

	_, err := tls.LoadX509KeyPair(param.certFile, param.keyFile)
	if err != nil {
		log.Errorf("Failed to load key pair: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(gin.Recovery())

	v1Group := router.Group("api/v1")
	{
		v1Group.POST("/mutate", handler.WebhookCallBack)
	}
	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 30 * time.Second
	endless.DefaultMaxHeaderBytes = 1 << 20

	srvAddr := ":" + strconv.Itoa(param.port)
	srv := endless.NewServer(srvAddr, router)

	err = srv.ListenAndServeTLS(param.certFile, param.keyFile)
	if err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Error("Listen: %s\n", err)
		}
	}
}
