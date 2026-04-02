package main

import (
	"github.com/gin-gonic/gin"

	"go-service/internal/handlers"
	"go-service/internal/middleware"
	"go-service/internal/store"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	s := store.NewStore()
	h := handlers.NewHandler(s)

	r.GET("/todos", h.GetTodos)
	r.POST("/todos", h.CreateTodo)
	r.GET("/todos/:id", h.GetTodoByID)
	r.PUT("/todos/:id", h.UpdateTodo)
	r.DELETE("/todos/:id", h.DeleteTodo)

	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
