package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"go-service/store"
)

type Handler struct {
	store *store.Store
}

type createTodoRequest struct {
	Title string `json:"title"`
}

type updateTodoRequest struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func NewHandler(store *store.Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) GetTodos(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetAll())
}

func (h *Handler) CreateTodo(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	todo := h.store.Create(req.Title)
	c.JSON(http.StatusCreated, todo)
}

func (h *Handler) GetTodoByID(c *gin.Context) {
	id := c.Param("id")
	todo, ok := h.store.GetByID(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *Handler) UpdateTodo(c *gin.Context) {
	id := c.Param("id")

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	updated, ok := h.store.Update(id, req.Title, req.Completed)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *Handler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	if !h.store.Delete(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
