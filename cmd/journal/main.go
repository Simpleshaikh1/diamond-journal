package main

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/Simpleshaikh1/diamond-journal/internal/server"
	"os"
	"path/filepath"

	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "Personal Journal CLI written in Go",
	Long:  `A beautiful, private, and extensible journaling tool.`,
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	//cli.Setup(rootCmd, cfg)

	// Setup repository (SQLite)
	dbPath := filepath.Join(cfg.StorageDir, cfg.DBFile)
	repo, err := repository.NewSQLiteRepository(dbPath)
	if err != nil {
		panic(err)
	}

	// Start API Server
	api := server.SetupAPI(repo)
	fmt.Println("🚀 Diamond Journal API running on http://localhost:8080")
	api.Run(":8080")

	//if err := rootCmd.Execute(); err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
}
