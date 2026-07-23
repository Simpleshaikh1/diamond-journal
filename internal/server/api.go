package server

import (
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/handler/http"
	"github.com/gin-gonic/gin"
)

func SetupAPI(repo domain.EntryRepository) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	h := http.NewEntryHandler(repo)

	api := r.Group("/api/v1")
	{
		entries := api.Group("/entries")
		{
			entries.GET("", h.List)
			entries.POST("", h.Create)
			entries.GET("/:id", h.GetByID)
			entries.PUT("/:id", h.Update)
			entries.DELETE("/:id", h.Delete)
			entries.GET("/search", h.Search)
		}
	}

	return r
}
