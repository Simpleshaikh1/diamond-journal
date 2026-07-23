package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newViewCmd(repo domain.EntryRepository) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view [id]",
		Short: "View a specific entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid ID")
				return
			}

			entry, err := repo.GetByID(uint(id))
			if err != nil {
				fmt.Printf("Entry not found: %v\n", err)
				return
			}

			fmt.Printf("\n#%d - %s\n", entry.ID, entry.Title)
			fmt.Printf("Date: %s | Mood: %s\n",
				entry.CreatedAt.Format("Monday, January 2, 2006 at 15:04"),
				entry.Mood)

			if len(entry.Tags) > 0 {
				fmt.Printf("Tags: %s\n", strings.Join(entry.Tags, ", "))
			}
			fmt.Println(strings.Repeat("─", 60))
			fmt.Println(entry.Content)
			fmt.Println(strings.Repeat("─", 60))
		},
	}
	return cmd
}
