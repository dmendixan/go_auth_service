package handlers

import (
	"auth-service/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

var adminTestDB *gorm.DB

func setupAdminTestEnv() *gin.Engine {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	adminTestDB = db

	_ = adminTestDB.AutoMigrate(&models.User{})

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/admin/users", func(c *gin.Context) { GetAllUsers(c, adminTestDB) })
	router.DELETE("/admin/users/:id", func(c *gin.Context) { DeleteUser(c, adminTestDB) })
	router.PUT("/admin/users/:id", func(c *gin.Context) { UpdateUser(c, adminTestDB) })

	return router
}

func TestGetAllUsersHandler(t *testing.T) {
	router := setupAdminTestEnv()

	adminTestDB.Create(&models.User{Name: "User1", Email: "user1@example.com", Password: "pass", Role: "user"})
	adminTestDB.Create(&models.User{Name: "User2", Email: "user2@example.com", Password: "pass", Role: "admin"})

	req, _ := http.NewRequest("GET", "/admin/users", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User1")
	assert.Contains(t, resp.Body.String(), "User2")
}

func TestDeleteUserHandler(t *testing.T) {
	router := setupAdminTestEnv()

	user := models.User{Name: "DeleteMe", Email: "deleteme@example.com", Password: "pass", Role: "user"}
	adminTestDB.Create(&user)

	req, _ := http.NewRequest("DELETE", "/admin/users/"+fmt.Sprint(user.ID), nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User deleted successfully")

	var deleted models.User
	result := adminTestDB.First(&deleted, user.ID)
	assert.Error(t, result.Error) // не найден
}

func TestUpdateUserHandler(t *testing.T) {
	router := setupAdminTestEnv()

	user := models.User{Name: "Old Name", Email: "old@example.com", Password: "pass", Role: "user"}
	adminTestDB.Create(&user)

	updateData := map[string]string{
		"name": "New Name",
		"role": "admin",
	}
	jsonValue, _ := json.Marshal(updateData)

	req, _ := http.NewRequest("PUT", "/admin/users/"+fmt.Sprint(user.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "User updated successfully")

	var updated models.User
	adminTestDB.First(&updated, user.ID)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "admin", updated.Role)
}
