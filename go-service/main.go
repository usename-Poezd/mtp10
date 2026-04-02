package main

import (
	"github.com/gin-gonic/gin"

	"go-service/handlers"
	"go-service/store"
)

func main() {
	r := gin.Default()

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
