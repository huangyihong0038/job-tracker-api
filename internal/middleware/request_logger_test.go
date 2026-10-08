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

func TestRequestLogger(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var output bytes.Buffer

	logger := log.New(
		&output,
		"",
		0,
	)

	router := gin.New()

	router.Use(
		RequestLogger(logger),
	)

	router.GET(
		"/test",
		func(c *gin.Context) {
			c.Status(http.StatusCreated)
		},
	)

	request := httptest.NewRequest(
		http.MethodGet,
		"/test?company=ByteDance",
		nil,
	)

	request.RemoteAddr =
		"192.0.2.1:12345"

	response := httptest.NewRecorder()

	router.ServeHTTP(
		response,
		request,
	)

	if response.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			response.Code,
		)
	}

	logLine := output.String()

	expectedParts := []string{
		"method=GET",
		`path="/test"`,
		"status=201",
		"latency=",
		"client_ip=192.0.2.1",
	}

	for _, expected := range expectedParts {
		if !strings.Contains(
			logLine,
			expected,
		) {
			t.Fatalf(
				"expected log to contain %q, got %q",
				expected,
				logLine,
			)
		}
	}
}
