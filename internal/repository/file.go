package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
)

type FileRepository struct {
	filepath string
	mu       sync.RWMutex
	entries  map[uint]*domain.Entry
	nextID   uint
}

func NewFileRepositoy(storageDir string) (*FileRepository, error) {
	filePath := filepath.Join(storageDir, "entries.json")
	repo := &FileRepository{
		filepath: filePath,
		entries:  make(map[uint]*domain.Entry),
		nextID:   1,
	}

	//Load existing data
	data, err := os.ReadFile(filePath)
	if err == nil {
		var list []*domain.Entry
		if jsonErr := json.Unmarshal(data, &list); jsonErr == nil {
			for _, e := range list {
				repo.entries[e.ID] = e
				if e.ID >= repo.nextID {
					repo.nextID = e.ID + 1
				}
			}
		}
	}

	return repo, nil
}
