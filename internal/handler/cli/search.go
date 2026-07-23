package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"strings"

	//"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/spf13/cobra"
)

func newSearchCmd(repo domain.EntryRepository) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search entries by keyword",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := strings.Join(args, " ")

			results, err := repo.Search(query)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if len(results) == 0 {
				fmt.Printf("No entries found for: \"%s\"\n", query)
				return
			}

			fmt.Printf("\n🔍 Search results for \"%s\" (%d found)\n", query, len(results))
			fmt.Println(strings.Repeat("─", 60))
			for _, e := range results {
				moodEmoji := getMoodEmoji(e.Mood)
				date := e.CreatedAt.Format("2006-01-02")

				fmt.Printf("%s %d. %s %s\n", moodEmoji, e.ID, date, e.Title)

				// Show snippet with highlight (simple version)
				snippet := e.Content
				if len(snippet) > 120 {
					snippet = snippet[:117] + "..."
				}
				fmt.Printf("   %s\n\n", snippet)
			}
		},
	}
	return cmd
}
