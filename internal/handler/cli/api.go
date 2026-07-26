package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/Simpleshaikh1/diamond-journal/internal/server"
	"github.com/spf13/cobra"
)

func NewAPICmd(repo domain.EntryRepository) *cobra.Command {
	var port string

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the REST API server",
		Run: func(cmd *cobra.Command, args []string) {
			if port == "" {
				port = "8080"
			}

			sqliteRepo, ok := repo.(*repository.SQLiteRepository)
			if !ok {
				fmt.Println("Error: Repository is not SQLite")
				return
			}

			api := server.SetupAPI(repo, sqliteRepo.GetDB()) // We'll adjust this

			fmt.Printf("🚀 Diamond Journal API started on http://localhost:%s\n", port)
			if err := api.Run(":" + port); err != nil {
				fmt.Printf("Server error: %v\n", err)
			}
		},
	}

	cmd.Flags().StringVarP(&port, "port", "p", "8080", "Port")
	return cmd
}
