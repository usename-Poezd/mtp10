package ws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func setupWSServer(t *testing.T) (*httptest.Server, *Hub) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	hub := NewHub()
	h := NewHandler(hub)
	r := gin.New()
	h.Init(r)
	srv := httptest.NewServer(r)
	return srv, hub
}

func dialWS(t *testing.T, srv *httptest.Server) *websocket.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial failed: %v", err)
	}
	return conn
}

func TestHubBroadcast(t *testing.T) {
	srv, _ := setupWSServer(t)
	defer srv.Close()

	c1 := dialWS(t, srv)
	defer c1.Close()
	c2 := dialWS(t, srv)
	defer c2.Close()

	msg := `{"type":"message","username":"Alice","text":"Hello"}`
	if err := c1.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	for _, conn := range []*websocket.Conn{c1, c2} {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		_, data, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if !strings.Contains(string(data), "Hello") {
			t.Errorf("expected 'Hello' in broadcast, got: %s", data)
		}
	}
}

func TestHubHealthEndpoint(t *testing.T) {
	srv, _ := setupWSServer(t)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ws")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestClientDisconnect(t *testing.T) {
	srv, hub := setupWSServer(t)
	defer srv.Close()

	c1 := dialWS(t, srv)
	c2 := dialWS(t, srv)

	c1.Close()
	time.Sleep(100 * time.Millisecond)

	msg := `{"type":"message","username":"Bob","text":"Still here"}`
	if err := c2.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("write after disconnect failed: %v", err)
	}
	c2.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := c2.ReadMessage()
	if err != nil {
		t.Fatalf("read failed after disconnect: %v", err)
	}
	if !strings.Contains(string(data), "Still here") {
		t.Errorf("expected message, got: %s", data)
	}
	_ = hub
	c2.Close()
}
