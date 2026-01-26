package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/velocitykode/velocity/pkg/router"
)

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{
		Request:  req,
		Response: rec,
	}

	err := Health(ctx)
	if err != nil {
		t.Fatalf("Health() returned error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Health() status = %d, want %d", rec.Code, http.StatusOK)
	}

	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Health() status = %q, want %q", response["status"], "healthy")
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	tests := []struct {
		name    string
		body    map[string]string
		wantErr string
	}{
		{
			name:    "missing all fields",
			body:    map[string]string{},
			wantErr: "Name, email, and password are required",
		},
		{
			name:    "missing email",
			body:    map[string]string{"name": "Test", "password": "secret"},
			wantErr: "Name, email, and password are required",
		},
		{
			name:    "missing password",
			body:    map[string]string{"name": "Test", "email": "test@example.com"},
			wantErr: "Name, email, and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := &router.Context{
				Request:  req,
				Response: rec,
			}

			_ = CreateUser(ctx)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("CreateUser() status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
			}

			var response map[string]string
			if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response["error"] != tt.wantErr {
				t.Errorf("CreateUser() error = %q, want %q", response["error"], tt.wantErr)
			}
		})
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	ctx := &router.Context{
		Request:  req,
		Response: rec,
	}

	_ = CreateUser(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("CreateUser() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestGetUser_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users/invalid", nil)
	rec := httptest.NewRecorder()
	ctx := &router.Context{
		Request:  req,
		Response: rec,
		Params:   map[string]string{"id": "invalid"},
	}

	_ = GetUser(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("GetUser() status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var response map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["error"] != "Invalid user ID" {
		t.Errorf("GetUser() error = %q, want %q", response["error"], "Invalid user ID")
	}
}
