package main

import (
	"github.com/gin-gonic/gin"

	"go-service/internal/gateway"
	"go-service/internal/middleware"
)

const gatewayPort = ":9000"

func main() {
	h := gateway.NewHandler()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	h.Init(r)

	if err := r.Run(gatewayPort); err != nil {
		panic(err)
	}
}
