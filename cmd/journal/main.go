package main

import (
	"fmt"
	"os"

	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/Simpleshaikh1/diamond-journal/internal/handler/cli"

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

	cli.Setup(rootCmd, cfg)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
