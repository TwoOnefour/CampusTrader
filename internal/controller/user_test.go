package controller

import (
	"CampusTrader/internal/common/database"
	"CampusTrader/internal/middleware/auth"
	"CampusTrader/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin(t *testing.T) {
	database.InitMySQL()
	userSrv := service.NewUserService(database.DB)
	userCtrl := NewUserController(userSrv)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/login", userCtrl.Login)

	body := `{"account": "testuser", "password": "password123"}`
	req, _ := http.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// 4. 创建响应记录器
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	t.Logf("Response: %s", w.Body.String())
}

func TestRegister(t *testing.T) {
	database.InitMySQL()
	userSrv := service.NewUserService(database.DB)
	userCtrl := NewUserController(userSrv)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", userCtrl.Register)
	body := `{"username": "testuser", "password": "password123", "re_password": "password123", "email": "testuser@gmail.com" }`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	t.Logf("Response: %s", w.Body.String())
}

func TestMe(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Error(".env file not found")
	}
	database.InitMySQL()
	userSrv := service.NewUserService(database.DB)
	userCtrl := NewUserController(userSrv)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(auth.JWTAuthMiddleware()).GET("/me", userCtrl.Me)
	req, _ := http.NewRequest("GET", "/me", nil)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiaXNzIjoibXktYXBwIiwiZXhwIjoxNzY1MTA4NTMxfQ.xcVTuZdB6yuccunOZlZrCRjRa75Au3sf-LzBW5bgTHc"
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `Bearer `+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	t.Logf("Response: %s", w.Body.String())
}
