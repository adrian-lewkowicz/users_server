package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Age   int    `json:"age" binding:"required"`
}

type Response struct {
	status   int
	jsonResp any
}

func setupDatabase() (*gorm.DB, error) {
	dsn := "host=localhost user=gorm password=gorm dbname=users port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	db, _ := setupDatabase()
	router := setUpRouter(db)
	router.Run(":8080")
}

var db *gorm.DB

func setUpRouter(database *gorm.DB) *gin.Engine {
	router := gin.Default()
	db = database
	db.AutoMigrate(&User{})

	var wg sync.WaitGroup

	router.GET("/users/:id", func(ctx *gin.Context) {
		handlehttp(ctx, &wg, getUserById)
	})

	router.POST("/users", func(ctx *gin.Context) {
		handlehttp(ctx, &wg, postUser)
	})

	router.PUT("/users/:id", func(ctx *gin.Context) {
		handlehttp(ctx, &wg, putUser)
	})

	router.DELETE("/users/:id", func(ctx *gin.Context) {
		handlehttp(ctx, &wg, deleteUser)
	})

	return router
}

func handlehttp(ctx *gin.Context, wg *sync.WaitGroup,
	f func(wg *sync.WaitGroup, ctx *gin.Context, ch chan Response)) {
	wg.Add(1)
	ch := make(chan Response)
	start := time.Now()
	go f(wg, ctx, ch)
	msg := <-ch
	close(ch)
	duration := time.Since(start)
	fmt.Printf("Time processing request: %v\n", duration)
	ctx.JSON(msg.status, msg.jsonResp)
}

func getUserById(wg *sync.WaitGroup, ctx *gin.Context, ch chan Response) {
	defer wg.Done()
	userID := ctx.Param("id")

	var user User
	result := db.Find(&user, userID)

	if result.Error != nil {
		ch <- Response{http.StatusNotFound, gin.H{"error": "Item not found"}}
		return
	} else {
		ch <- Response{http.StatusOK, user}
		return
	}
}

func postUser(wg *sync.WaitGroup, ctx *gin.Context, ch chan Response) {
	defer wg.Done()
	var user User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ch <- Response{http.StatusBadRequest, gin.H{"error": err.Error()}}
		return
	}

	result := db.Create(&user)
	if result.Error != nil {
		ch <- Response{http.StatusInternalServerError, gin.H{"error": "Failed to create resource"}}
		return
	}
	ch <- Response{http.StatusOK, user}
}

func putUser(wg *sync.WaitGroup, ctx *gin.Context, ch chan Response) {
	defer wg.Done()
	userID := ctx.Param("id")
	var updatedUser User

	if err := ctx.ShouldBindBodyWithJSON(&updatedUser); err != nil {
		ch <- Response{http.StatusBadRequest, gin.H{"error": err.Error()}}
		return
	}

	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		ch <- Response{http.StatusNotFound, gin.H{"error": "Item not found"}}
		return
	}

	result = db.Model(&User{}).Where("id = ?", userID).Updates(updatedUser)
	if result.Error != nil {
		ch <- Response{http.StatusInternalServerError, gin.H{"error": "Failed to update resource"}}
		return
	}

	ch <- Response{http.StatusOK, updatedUser}
}

func deleteUser(wg *sync.WaitGroup, ctx *gin.Context, ch chan Response) {
	defer wg.Done()
	userID := ctx.Param("id")

	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		ch <- Response{http.StatusNotFound, gin.H{"error": "Item not found"}}
		return
	}
	result = db.Delete(user)
	if result.Error != nil {
		ch <- Response{http.StatusInternalServerError, gin.H{"error": "Failed to delete"}}
		return
	}
	ch <- Response{http.StatusOK, "ok"}
}
