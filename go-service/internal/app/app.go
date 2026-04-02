package app

import (
	"github.com/gin-gonic/gin"

	handler "go-service/internal/handlers/http"
	"go-service/internal/middleware"
	"go-service/internal/store"
)

func Run() {
	s := store.NewStore()
	h := handler.NewHandler(s)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	h.Init(r)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
