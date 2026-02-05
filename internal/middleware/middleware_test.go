package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/velocitykode/velocity/pkg/router"
)

func TestCORSMiddleware(t *testing.T) {
	handler := CORSMiddleware(func(c *router.Context) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("CORSMiddleware() returned error: %v", err)
	}

	tests := []struct {
		header string
		want   string
	}{
		{"Access-Control-Allow-Origin", "*"},
		{"Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS"},
		{"Access-Control-Allow-Credentials", "true"},
	}

	for _, tt := range tests {
		got := rec.Header().Get(tt.header)
		if got != tt.want {
			t.Errorf("CORSMiddleware() %s = %q, want %q", tt.header, got, tt.want)
		}
	}
}

func TestCORSMiddleware_OPTIONS(t *testing.T) {
	nextCalled := false
	handler := CORSMiddleware(func(c *router.Context) error {
		nextCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("CORSMiddleware() returned error: %v", err)
	}

	if nextCalled {
		t.Error("CORSMiddleware() called next handler for OPTIONS request")
	}

	if rec.Code != http.StatusOK {
		t.Errorf("CORSMiddleware() status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestEnsureJSONMiddleware(t *testing.T) {
	handler := EnsureJSONMiddleware(func(c *router.Context) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("EnsureJSONMiddleware() returned error: %v", err)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("EnsureJSONMiddleware() Content-Type = %q, want %q", contentType, "application/json")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	nextCalled := false
	handler := LoggingMiddleware(func(c *router.Context) error {
		nextCalled = true
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("LoggingMiddleware() returned error: %v", err)
	}

	if !nextCalled {
		t.Error("LoggingMiddleware() did not call next handler")
	}
}


func TestTrustProxiesMiddleware(t *testing.T) {
	handler := TrustProxiesMiddleware(func(c *router.Context) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1")
	req.Header.Set("X-Forwarded-Proto", "https")
	req.Header.Set("X-Forwarded-Host", "example.com")
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("TrustProxiesMiddleware() returned error: %v", err)
	}

	if req.RemoteAddr != "192.168.1.1" {
		t.Errorf("TrustProxiesMiddleware() RemoteAddr = %q, want %q", req.RemoteAddr, "192.168.1.1")
	}

	if req.URL.Scheme != "https" {
		t.Errorf("TrustProxiesMiddleware() Scheme = %q, want %q", req.URL.Scheme, "https")
	}

	if req.Host != "example.com" {
		t.Errorf("TrustProxiesMiddleware() Host = %q, want %q", req.Host, "example.com")
	}
}

func TestPreventRequestsDuringMaintenanceMiddleware(t *testing.T) {
	t.Run("no maintenance", func(t *testing.T) {
		nextCalled := false
		handler := PreventRequestsDuringMaintenanceMiddleware(func(c *router.Context) error {
			nextCalled = true
			return nil
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := &router.Context{Request: req, Response: rec}

		if err := handler(ctx); err != nil {
			t.Fatalf("PreventRequestsDuringMaintenanceMiddleware() returned error: %v", err)
		}

		if !nextCalled {
			t.Error("PreventRequestsDuringMaintenanceMiddleware() did not call next handler")
		}
	})

	t.Run("in maintenance", func(t *testing.T) {
		// Create maintenance file
		dir := filepath.Join("storage", "framework")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create maintenance directory: %v", err)
		}
		defer os.RemoveAll("storage")

		file := filepath.Join(dir, "down")
		if err := os.WriteFile(file, []byte{}, 0644); err != nil {
			t.Fatalf("Failed to create maintenance file: %v", err)
		}

		nextCalled := false
		handler := PreventRequestsDuringMaintenanceMiddleware(func(c *router.Context) error {
			nextCalled = true
			return nil
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := &router.Context{Request: req, Response: rec}

		if err := handler(ctx); err != nil {
			t.Fatalf("PreventRequestsDuringMaintenanceMiddleware() returned error: %v", err)
		}

		if nextCalled {
			t.Error("PreventRequestsDuringMaintenanceMiddleware() called next handler during maintenance")
		}

		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("PreventRequestsDuringMaintenanceMiddleware() status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
	})
}

func TestValidatePostSizeMiddleware(t *testing.T) {
	handler := ValidatePostSizeMiddleware(100)(func(c *router.Context) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	if err := handler(ctx); err != nil {
		t.Fatalf("ValidatePostSizeMiddleware() returned error: %v", err)
	}
}

func TestAuth_Unauthorized(t *testing.T) {
	handler := Auth(func(c *router.Context) error {
		return nil
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{Request: req, Response: rec}

	_ = handler(ctx)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("Auth() status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}
