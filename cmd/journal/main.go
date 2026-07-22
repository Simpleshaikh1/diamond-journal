package main

import (
	"fmt"
	"os"

	"github.com/Simpleshaikh1/diamond-journal/handler/cli"
	"github.com/Simpleshaikh1/diamond-journal/internal/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "journal",
	Short: "Personal Journal CLI",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
