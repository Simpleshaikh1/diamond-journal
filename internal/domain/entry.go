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
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Mood     Mood     `json:"mood"`
	Tags     []string `json:"tags" gorm:"serialize:json"`
	Location string   `json:"location,omitempty"`
}

type EntryRepository interface {
	Create(entry *Entry) error
	GetByID(id uint) (*Entry, error)
	List(page, limit int, filters map[string]interface{}) ([]Entry, int64, error)
	Update(entry *Entry) error
	Delete(id uint) error
	Search(query string) ([]Entry, error)
}
