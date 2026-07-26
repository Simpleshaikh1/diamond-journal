package http

import (
	"net/http"
	"strconv"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/gin-gonic/gin"
)

type EntryHandler struct {
	repo domain.EntryRepository
}

func NewEntryHandler(repo domain.EntryRepository) *EntryHandler {
	return &EntryHandler{repo: repo}
}

// Get all entries
func (h *EntryHandler) List(c *gin.Context) {
	var filter domain.EntryFilter

	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	entries, total, err := h.repo.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"entries": entries,
			"pagination": gin.H{
				"page":       filter.Page,
				"limit":      filter.Limit,
				"total":      total,
				"totalPages": totalPages,
			},
		},
	})
}

// Create entry
func (h *EntryHandler) Create(c *gin.Context) {
	var entry domain.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Entry created", "data": entry})
}

// Get one entry
func (h *EntryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	entry, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": entry})
}

// Update Entry
func (h *EntryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var entry domain.Entry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.ID = uint(id)
	if err := h.repo.Update(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Entry updated",
		"data":    entry,
	})
}

// DELETE Entry
func (h *EntryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Entry deleted",
	})
}

// Search
func (h *EntryHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	entries, err := h.repo.Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   entries,
	})
}
