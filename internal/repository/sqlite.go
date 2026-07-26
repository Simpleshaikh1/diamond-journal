package repository

import (
	"fmt"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"

	//"gorm.io/driver/sqlite"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	sqliteDriver "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *gorm.DB
}

func (r *SQLiteRepository) GetDB() *gorm.DB {
	return r.db
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dsn := dbPath

	if !filepath.IsAbs(dbPath) {
		dsn = filepath.Join(".", dbPath)
	}

	// Use modernc.org/sqlite driver (pure Go)
	db, err := gorm.Open(sqliteDriver.Dialector{
		DriverName: "sqlite",
		DSN:        dsn,
		Conn:       nil, //gorm will handle connection
	}, &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&domain.Entry{}, &domain.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}
	fmt.Println("✅ Connected to SQLite database")
	return &SQLiteRepository{db: db}, nil
}

func (r *SQLiteRepository) Create(entry *domain.Entry) error {
	return r.db.Create(entry).Error
}

func (r *SQLiteRepository) GetByID(id uint) (*domain.Entry, error) {
	var entry domain.Entry
	err := r.db.First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *SQLiteRepository) List(filter domain.EntryFilter) ([]domain.Entry, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Page < 1 {
		filter.Page = 1
	}

	offset := (filter.Page - 1) * filter.Limit

	var entries []domain.Entry
	var total int64

	query := r.db.Model(&domain.Entry{}).Order("created_at desc")

	// Apply filters
	if filter.Mood != "" {
		query = query.Where("mood = ?", filter.Mood)
	}
	if filter.Tag != "" {
		query = query.Where("json_contains(tags, ?)", `["`+filter.Tag+`"]`)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", search, search)
	}

	// Count total
	query.Count(&total)

	// Get data
	err := query.Offset(offset).Limit(filter.Limit).Find(&entries).Error

	return entries, total, err
}

func (r *SQLiteRepository) Update(entry *domain.Entry) error {
	return r.db.Save(entry).Error
}

func (r *SQLiteRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Entry{}, id).Error
}

func (r *SQLiteRepository) Search(query string) ([]domain.Entry, error) {
	var entries []domain.Entry
	err := r.db.Where("title LIKE ? OR content LIKE ?",
		"%"+query+"%", "%"+query+"%").
		Order("created_at desc").
		Find(&entries).Error
	return entries, err
}
