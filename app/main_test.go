package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestGet(t *testing.T) {
	db, _ := setupMockDB()
	router := setUpRouter(db)

	t.Run("test post ok", func(t *testing.T) {
		w := httptest.NewRecorder()
		user := User{
			Name:  "John",
			Email: "email@email.com",
			Age:   33,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fail()
			return
		}

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("test post Bad Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		user := User{
			Name: "John",
			Age:  33,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fail()
			return
		}
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("test get by id", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("test put ok", func(t *testing.T) {

		w := httptest.NewRecorder()
		user := User{
			Name:  "John",
			Email: "email@email.com",
			Age:   33,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fail()
			return
		}
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(jsonData), w.Body.String())
	})

	t.Run("test put not found", func(t *testing.T) {

		w := httptest.NewRecorder()
		user := User{
			Name:  "John",
			Email: "email@email.com",
			Age:   33,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fail()
			return
		}
		req, _ := http.NewRequest("PUT", "/users/0", bytes.NewBuffer(jsonData))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("test put not found", func(t *testing.T) {

		w := httptest.NewRecorder()
		user := User{
			Email: "email@email.com",
			Age:   33,
		}
		jsonData, err := json.Marshal(user)
		if err != nil {
			t.Fail()
			return
		}
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("test DELETE", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("test DELETE not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/users/0", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func setupMockDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&User{})
	return db, err
}
