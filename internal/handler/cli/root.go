package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

func Setup(root *cobra.Command, cfg *config.Config) {
	var repo domain.EntryRepository
	var err error

	switch strings.ToLower(cfg.StorageType) {
	case "file", "json":
		repo, err = repository.NewFileRepository(cfg.StorageDir)
		fmt.Println("📁 Using File (JSON) storage")

	case "sqlite", "db", "":
		// Default to SQLite
		dbPath := filepath.Join(cfg.StorageDir, cfg.DBFile)
		repo, err = repository.NewSQLiteRepository(dbPath)
		fmt.Println("🗄️  Using SQLite database")

	default:
		fmt.Printf("Unknown storage_type: %s. Defaulting to SQLite.\n", cfg.StorageType)
		dbPath := filepath.Join(cfg.StorageDir, cfg.DBFile)
		repo, err = repository.NewSQLiteRepository(dbPath)
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to initialize repository: %v", err))
	}

	//Add commands
	root.AddCommand(newCreateCmd(repo))
	root.AddCommand(newListCmd(repo))
	root.AddCommand(newViewCmd(repo))
	root.AddCommand(newEditCmd(repo))
	root.AddCommand(newDeleteCmd(repo))
	root.AddCommand(newSearchCmd(repo))
	root.AddCommand(newAPICmd(repo))
}
