package handlers

import (
	"auth-service/config"
	"auth-service/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite" // 👈 используем вместо gorm.io/driver/sqlite
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&models.User{}, &models.RefreshToken{}) // ✅ теперь обе таблицы
	return db
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	router := gin.Default()

	// Используем обёртку RegisterWithDB
	router.POST("/register", RegisterWithDB(db))

	payload := `{
		"name": "Test User",
		"email": "test@example.com",
		"password": "123456"
	}`

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "User registered successfully", response["message"])
}
func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db := setupTestDB()
	router := gin.Default()

	// Регистрируем пользователя напрямую в БД
	password := "mypassword"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Name:         "Test User",
		Email:        "login@example.com",
		PasswordHash: string(hashed),
		Role:         "user",
	}
	db.Create(&user)

	// Регистрируем endpoint логина
	router.POST("/login", func(c *gin.Context) {
		configBackup := config.DB // если ты всё ещё используешь config.DB в Login
		config.DB = db
		defer func() { config.DB = configBackup }()

		Login(c)
	})

	// Формируем запрос логина
	payload := `{
		"email": "login@example.com",
		"password": "mypassword"
	}`

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Contains(t, resp, "access_token")
	assert.Contains(t, resp, "refresh_token")
}
