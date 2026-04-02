package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLoggerWritesStructuredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	t.Cleanup(func() {
		log.SetOutput(originalOutput)
	})

	r := gin.New()
	r.Use(Logger())
	r.GET("/todos", func(c *gin.Context) {
		c.JSON(http.StatusOK, []string{})
	})

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	output := buf.String()
	checks := []string{"method=GET", "path=/todos", "status=200", "latency_ms=", "ip="}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Fatalf("expected log output to contain %q, got %q", check, output)
		}
	}
}
