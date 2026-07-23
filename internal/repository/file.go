package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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

func (r *FileRepository) Delete(id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[id]; !exists {
		return fmt.Errorf("entry not found")
	}

	delete(r.entries, id)
	return r.save()
}

func (r *FileRepository) GetByID(id uint) (*domain.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if entry, exists := r.entries[id]; exists {
		return entry, nil
	}
	return nil, fmt.Errorf("entry not found")
}

func (r *FileRepository) List(page, limit int, _ map[string]interface{}) ([]domain.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if page < 1 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	var entries []domain.Entry
	for _, e := range r.entries {
		entries = append(entries, *e)
	}

	// Simple sort by ID descending (newest first)
	//for i := 0; i < len(entries)-1; i++ {
	//	for j := i + 1; j < len(entries); j++ {
	//		if entries[i].ID < entries[j].ID {
	//			entries[i], entries[j] = entries[j], entries[i]
	//		}
	//	}
	//}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID > entries[i].ID
	})

	// Calculate offset
	offset := (page - 1) * limit

	// Apply pagination
	if offset >= len(entries) {
		return []domain.Entry{}, nil
	}

	end := offset + limit
	if end > len(entries) {
		end = len(entries)
	}

	//if len(entries) > limit {
	//	entries = entries[:limit]
	//}

	//return entries, nil

	return entries[offset:end], nil
}

func (r *FileRepository) Search(query string) ([]domain.Entry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query = strings.ToLower(strings.TrimSpace(query))
	var results []domain.Entry

	for _, entry := range r.entries {
		contentLower := strings.ToLower(entry.Content)
		titleLower := strings.ToLower(entry.Title)

		if strings.Contains(titleLower, query) || strings.Contains(contentLower, query) {
			results = append(results, *entry)
		}
	}

	// Sort newest first
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].CreatedAt.Before(results[j].CreatedAt) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

func (r *FileRepository) save() error {
	list := make([]*domain.Entry, 0, len(r.entries))
	for _, e := range r.entries {
		list = append(list, e)
	}
	data, _ := json.MarshalIndent(list, "", " ")
	return os.WriteFile(r.filepath, data, 0644)
}
