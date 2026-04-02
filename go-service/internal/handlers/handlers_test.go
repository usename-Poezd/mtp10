package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"go-service/internal/store"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	s := store.NewStore()
	h := NewHandler(s)
	r := gin.New()

	r.GET("/todos", h.GetTodos)
	r.POST("/todos", h.CreateTodo)
	r.GET("/todos/:id", h.GetTodoByID)
	r.PUT("/todos/:id", h.UpdateTodo)
	r.DELETE("/todos/:id", h.DeleteTodo)

	return r
}

func TestGetTodosReturnsEmptyArray(t *testing.T) {
	r := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if w.Body.String() != "[]" {
		t.Fatalf("expected [] body, got %s", w.Body.String())
	}
}

func TestCreateAndGetTodo(t *testing.T) {
	r := setupRouter()
	body := []byte(`{"title":"Test Task"}`)

	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	if createW.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, createW.Code)
	}

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	id, ok := created["id"].(string)
	if !ok || id == "" {
		t.Fatal("expected non-empty id")
	}

	getReq := httptest.NewRequest(http.MethodGet, "/todos/"+id, nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getW.Code)
	}
}

func TestCreateTodoValidation(t *testing.T) {
	r := setupRouter()
	body := []byte(`{"title":""}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateAndDeleteTodo(t *testing.T) {
	r := setupRouter()
	body := []byte(`{"title":"Task"}`)

	createReq := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(body))
	createReq.Header.Set("Content-Type", "application/json")
	createW := httptest.NewRecorder()
	r.ServeHTTP(createW, createReq)

	var created map[string]interface{}
	if err := json.Unmarshal(createW.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to parse create response: %v", err)
	}
	id := created["id"].(string)

	updateBody := []byte(`{"title":"Updated","completed":true}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/"+id, bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)

	if updateW.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, updateW.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/"+id, nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)

	if deleteW.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, deleteW.Code)
	}
}

func TestNotFoundCases(t *testing.T) {
	r := setupRouter()

	getReq := httptest.NewRequest(http.MethodGet, "/todos/missing", nil)
	getW := httptest.NewRecorder()
	r.ServeHTTP(getW, getReq)
	if getW.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, getW.Code)
	}

	updateBody := []byte(`{"title":"Updated","completed":true}`)
	updateReq := httptest.NewRequest(http.MethodPut, "/todos/missing", bytes.NewBuffer(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateW := httptest.NewRecorder()
	r.ServeHTTP(updateW, updateReq)
	if updateW.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, updateW.Code)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/todos/missing", nil)
	deleteW := httptest.NewRecorder()
	r.ServeHTTP(deleteW, deleteReq)
	if deleteW.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, deleteW.Code)
	}
}
