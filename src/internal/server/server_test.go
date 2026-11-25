package server

import (
	"avan/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestNewServerDoesNotTrustForwardedClientIPByDefault(t *testing.T) {
	server := NewServer(config.Config{}, zap.Logger{})
	server.Engine.GET("/client-ip", func(c *gin.Context) {
		c.String(http.StatusOK, c.ClientIP())
	})
	request := httptest.NewRequest(http.MethodGet, "/client-ip", nil)
	request.RemoteAddr = "198.51.100.7:4567"
	request.Header.Set("X-Forwarded-For", "203.0.113.9")
	response := httptest.NewRecorder()
	server.Engine.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Body.String() != "198.51.100.7" {
		t.Fatalf("response = %d %q, want direct peer IP", response.Code, response.Body.String())
	}
}
