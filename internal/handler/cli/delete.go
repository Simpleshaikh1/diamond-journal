package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func newDeleteCmd(repo *repository.FileRepository) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete an entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("❌ Invalid ID")
				return
			}

			// Confirm
			fmt.Printf("Delete entry #%d? (y/N): ", id)
			var confirm string
			fmt.Scanln(&confirm)

			if strings.ToLower(confirm) != "y" {
				fmt.Println("Cancelled.")
				return
			}

			if err := repo.Delete(uint(id)); err != nil {
				fmt.Printf("❌ %v\n", err)
				return
			}

			fmt.Printf("🗑️  Entry #%d deleted.\n", id)
		},
	}
	return cmd
}
