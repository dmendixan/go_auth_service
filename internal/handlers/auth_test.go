package handlers

import (
	"auth-service/internal/models"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Временная переменная для базы
var testDB *gorm.DB

// Инициализация базы и роутера перед тестами
func setupTestEnv() *gin.Engine {
	// Используем SQLite в памяти
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	testDB = db

	// Миграция модели пользователя
	err = testDB.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate user model")
	}

	// Подключаем gin
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Регистрируем роут (только для теста)
	router.POST("/register", func(c *gin.Context) {
		Register(c, testDB) // передаем временную DB
	})

	return router
}

// Тест для /register
func TestRegisterHandler(t *testing.T) {
	router := setupTestEnv()

	// Тестовые данные
	body := map[string]string{
		"email":    "test@example.com",
		"password": "123456",
		"name":     "Test User",
	}
	jsonValue, _ := json.Marshal(body)

	// Создаем запрос
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Отправляем запрос
	router.ServeHTTP(resp, req)

	// Проверки
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Contains(t, resp.Body.String(), "User registered successfully")
}
