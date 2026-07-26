package server

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/handler/http"
	"github.com/Simpleshaikh1/diamond-journal/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupAPI(repo domain.EntryRepository, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	if db != nil {
		fmt.Println("Warning: DB is nil in SetupApi")
	}

	// Initialize Auth
	authHandler := http.NewAuthHandler(db)
	entryHandler := http.NewEntryHandler(repo)

	// Public routes (Auth)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.JWTAuth())
	{
		entries := protected.Group("/entries")
		{
			entries.GET("", entryHandler.List)
			entries.POST("", entryHandler.Create)
			entries.GET("/:id", entryHandler.GetByID)
			entries.PUT("/:id", entryHandler.Update)
			entries.DELETE("/:id", entryHandler.Delete)
			entries.GET("/search", entryHandler.Search)
		}
	}

	return r
}
