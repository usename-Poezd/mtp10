package main

import (
	"github.com/gin-gonic/gin"

	"go-service/internal/middleware"
	"go-service/internal/ws"
)

const wsPort = ":8081"

func main() {
	hub := ws.NewHub()
	h := ws.NewHandler(hub)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	h.Init(r)

	if err := r.Run(wsPort); err != nil {
		panic(err)
	}
}
