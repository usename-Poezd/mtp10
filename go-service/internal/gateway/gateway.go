package gateway

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	defaultGoBackendURL     = "http://localhost:8080"
	defaultPythonBackendURL = "http://localhost:8000"
)

type Handler struct {
	goBackendURL     string
	pythonBackendURL string
}

func NewHandler() *Handler {
	goURL := os.Getenv("GO_BACKEND_URL")
	if goURL == "" {
		goURL = defaultGoBackendURL
	}

	pythonURL := os.Getenv("PYTHON_BACKEND_URL")
	if pythonURL == "" {
		pythonURL = defaultPythonBackendURL
	}

	return &Handler{
		goBackendURL:     goURL,
		pythonBackendURL: pythonURL,
	}
}

func (h *Handler) Init(r *gin.Engine) {
	r.GET("/health", h.Health)
	r.Any("/api/go/*path", h.ProxyToGo)
	r.Any("/api/python/*path", h.ProxyToPython)
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) ProxyToGo(c *gin.Context) {
	path := c.Param("path")
	h.proxy(c, h.goBackendURL, path)
}

func (h *Handler) ProxyToPython(c *gin.Context) {
	path := c.Param("path")
	h.proxy(c, h.pythonBackendURL, path)
}

func (h *Handler) proxy(c *gin.Context, targetBase, path string) {
	target, err := url.Parse(targetBase)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid backend URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":  "bad gateway",
			"detail": "backend unavailable",
		})
	}

	c.Request.URL.Path = path
	c.Request.URL.Host = target.Host
	c.Request.URL.Scheme = target.Scheme
	c.Request.Host = target.Host

	proxy.ServeHTTP(c.Writer, c.Request)
}
