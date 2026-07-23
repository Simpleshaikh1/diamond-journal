package server

import (
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/handler/http"
	"github.com/gin-gonic/gin"
)

func SetupAPI(repo domain.EntryRepository) *gin.Engine {
	r := gin.Default()

	handler := http.NewEntryHandler(repo)

	api := r.Group("/api/v1")
	{
		api.GET("/entries", handler.List)
		api.POST("/entries", handler.Create)
		api.GET("/entries/:id", handler.GetByID)
		// We'll add more (update, delete, search) soon
	}

	return r
}
