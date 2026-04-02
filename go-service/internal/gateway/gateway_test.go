package gateway

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupGateway(goURL, pythonURL string) *gin.Engine {
	gin.SetMode(gin.TestMode)

	h := &Handler{
		goBackendURL:     goURL,
		pythonBackendURL: pythonURL,
	}

	r := gin.New()
	h.Init(r)

	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupGateway("http://localhost:8080", "http://localhost:8000")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status=ok, got %q", resp["status"])
	}
}

func TestProxyToGo(t *testing.T) {
	goBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todos" {
			t.Fatalf("expected path /todos, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"1","title":"go"}]`))
	}))
	defer goBackend.Close()

	r := setupGateway(goBackend.URL, "http://localhost:8000")
	gatewayServer := httptest.NewServer(r)
	defer gatewayServer.Close()

	resp, err := http.Get(gatewayServer.URL + "/api/go/todos")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if string(body) != `[{"id":"1","title":"go"}]` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestProxyToPython(t *testing.T) {
	pythonBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/todos" {
			t.Fatalf("expected path /todos, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"id":"2","title":"python"}]`))
	}))
	defer pythonBackend.Close()

	r := setupGateway("http://localhost:8080", pythonBackend.URL)
	gatewayServer := httptest.NewServer(r)
	defer gatewayServer.Close()

	resp, err := http.Get(gatewayServer.URL + "/api/python/todos")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	if string(body) != `[{"id":"2","title":"python"}]` {
		t.Fatalf("unexpected body: %s", string(body))
	}
}

func TestUnknownRoute(t *testing.T) {
	r := setupGateway("http://localhost:8080", "http://localhost:8000")

	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestBackendUnavailable(t *testing.T) {
	goBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	backendURL := goBackend.URL
	goBackend.Close()

	r := setupGateway(backendURL, "http://localhost:8000")
	gatewayServer := httptest.NewServer(r)
	defer gatewayServer.Close()

	resp, err := http.Get(gatewayServer.URL + "/api/go/todos")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected status %d, got %d", http.StatusBadGateway, resp.StatusCode)
	}
}
