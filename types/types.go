package types

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type ProtectNS []string

func (i *ProtectNS) String() string {
	return fmt.Sprint(*i)
}

func (i *ProtectNS) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type WHSrvParam struct {
	Port     int
	CertFile string
	KeyFile  string
	ProtectNS
}

type NewContext struct {
	*gin.Context
	ProtectNS []string `json:"protect_ns"`
}

type HandlerFunc func(ctx *NewContext)
