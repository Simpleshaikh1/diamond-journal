package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Simpleshaikh1/diamond-journal/internal/auth"
	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/Simpleshaikh1/diamond-journal/internal/handler/cli"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "Diamond Journal - CLI + API",
	Long:  `A beautiful, private, and extensible journaling tool.`,
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Initialize JWT
	auth.Init(cfg.JWTSecret)

	// Setup repository
	dbPath := filepath.Join(cfg.StorageDir, cfg.DBFile)
	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		panic(err)
	}

	// Setup CLI (this already adds most commands)
	cli.Setup(rootCmd, cfg)

	// Add API Command
	rootCmd.AddCommand(cli.NewAPICmd(repo))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
