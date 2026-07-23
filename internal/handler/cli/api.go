package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/server"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func newAPICmd(repo domain.EntryRepository) *cobra.Command {
	var port string

	cmd := &cobra.Command{
		Use:   "api",
		Short: "Start the REST API server",
		Run: func(cmd *cobra.Command, args []string) {
			if port == "" {
				port = "8080"
			}

			gin.SetMode(gin.DebugMode) // Change to ReleaseMode in production

			api := server.SetupAPI(repo)

			// Graceful shutdown
			go func() {
				fmt.Printf("🚀 Diamond Journal API started on http://localhost:%s\n", port)
				fmt.Println("Press Ctrl+C to stop")
				if err := api.Run(":" + port); err != nil {
					fmt.Printf("Server error: %v\n", err)
				}
			}()

			// Wait for interrupt signal
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			fmt.Println("\nShutting down server...")
		},
	}

	cmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run the API on")
	return cmd
}
