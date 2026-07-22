package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
)

type FileRepository struct {
	filepath string
	mu       sync.RWMutex
	entries  map[uint]*domain.Entry
	nextID   uint
}

func NewFileRepository(storageDir string) (*FileRepository, error) {
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

func (r *FileRepository) Create(entry *domain.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry.ID = r.nextID
	now := time.Now()
	entry.CreatedAt = now
	entry.UpdatedAt = now

	r.entries[entry.ID] = entry
	r.nextID++

	return r.save()
}

func (r *FileRepository) Update(entry *domain.Entry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	//check if entry exist
	existing, ok := r.entries[entry.ID]
	if !ok {
		return fmt.Errorf("entry with ID %d not found", entry.ID)
	}

	// Update fields
	existing.Title = entry.Title
	existing.Content = entry.Content
	existing.Mood = entry.Mood
	existing.Tags = entry.Tags
	existing.Location = entry.Location
	existing.UpdatedAt = time.Now()

	return r.save()
}

func (r *FileRepository) save() error {
	list := make([]*domain.Entry, 0, len(r.entries))
	for _, e := range r.entries {
		list = append(list, e)
	}
	data, _ := json.MarshalIndent(list, "", " ")
	return os.WriteFile(r.filepath, data, 0644)
}
