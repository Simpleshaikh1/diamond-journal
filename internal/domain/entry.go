package domain

import (
	"time"
)

type Mood string

const (
	MoodGreat    Mood = "great"
	MoodGood     Mood = "good"
	MoodOkay     Mood = "okay"
	MoodBad      Mood = "bad"
	MoodTerrible Mood = "terrible"
)

type Entry struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`

	Title    string   `json:"title" gorm:"not null"`
	Content  string   `json:"content"`
	Mood     Mood     `json:"mood"`
	Tags     []string `json:"tags" gorm:"serializer:json"`
	Location string   `json:"location,omitempty"`
}

type EntryFilter struct {
	Page      int
	Limit     int
	Mood      Mood
	Tag       string
	StartDate *time.Time
	EndDate   *time.Time
	Search    string
}
type EntryRepository interface {
	Create(entry *Entry) error
	GetByID(id uint) (*Entry, error)
	List(page, limit int, filters map[string]interface{}) ([]Entry, int64, error)
	Update(entry *Entry) error
	Delete(id uint) error
	Search(query string) ([]Entry, error)
}
