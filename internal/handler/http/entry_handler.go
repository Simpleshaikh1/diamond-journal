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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	entries, err := h.repo.List(1, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entries})
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

	c.JSON(http.StatusCreated, gin.H{"message": "Entry created", "data": entry})
}

// Get one entry
func (h *EntryHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	entry, err := h.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}
