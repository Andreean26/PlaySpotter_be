package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"playspotter/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Test registration flow
func TestRegisterFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create test router
	router := gin.New()

	// Mock auth handler (without DB for unit test)
	router.POST("/auth/register", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "validation_error", "message": err.Error()}})
			return
		}

		// Mock successful registration
		c.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"id":    uuid.New(),
				"name":  req.Name,
				"email": req.Email,
				"role":  "user",
			},
		})
	})

	// Test valid registration
	t.Run("Valid Registration", func(t *testing.T) {
		body := map[string]string{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["role"] != "user" {
			t.Errorf("Expected role 'user', got %s", data["role"])
		}
	})

	// Test invalid email
	t.Run("Invalid Email", func(t *testing.T) {
		body := map[string]string{
			"name":     "Test User",
			"email":    "invalid-email",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code == http.StatusCreated {
			t.Error("Expected validation error for invalid email")
		}
	})
}

// Test login flow
func TestLoginFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.POST("/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "validation_error", "message": err.Error()}})
			return
		}

		// Mock successful login
		if req.Email == "test@example.com" && req.Password == "password123" {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					"access_token":  "mock_access_token",
					"refresh_token": "mock_refresh_token",
					"user": gin.H{
						"id":    uuid.New(),
						"email": req.Email,
						"role":  "user",
					},
				},
			})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"code":    "invalid_credentials",
				"message": "Invalid credentials",
			},
		})
	})

	t.Run("Valid Login", func(t *testing.T) {
		body := map[string]string{
			"email":    "test@example.com",
			"password": "password123",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["access_token"] == "" {
			t.Error("Expected access_token in response")
		}
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		body := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})
}

// Test health check
func TestHealthCheck(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/health", handlers.HealthCheck)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	if data["ok"] != true {
		t.Error("Expected health check to return ok: true")
	}
}
