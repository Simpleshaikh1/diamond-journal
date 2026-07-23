package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
)

func newListCmd(repo *repository.FileRepository) *cobra.Command {
	var limit int
	var moodStr, tag string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Journal Entries",
		Run: func(cmd *cobra.Command, args []string) {
			entries, err := repo.List(1, limit, nil)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			if len(entries) == 0 {
				fmt.Println("No entries yet. Create one with `journal new`")
				return
			}

			fmt.Printf("\n📖 Your Journal Entries (%d total)\n", len(entries))
			fmt.Println(strings.Repeat("─", 60))

			for _, e := range entries {
				moodEmoji := getMoodEmoji(e.Mood)
				date := e.CreatedAt.Format("2006-01-02")
				tags := ""
				if len(e.Tags) > 0 {
					tags = " #" + strings.Join(e.Tags, " #")
				}

				fmt.Printf("%d. %s %s %s%s\n",
					e.ID, moodEmoji, date, e.Title, tags)

				// Preview of content
				preview := e.Content
				if len(preview) > 80 {
					preview = preview[:77] + "..."
				}
				if preview != "" {
					fmt.Printf("   %s\n", preview)
				}
				fmt.Println()
			}
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Number of entries to show")
	cmd.Flags().StringVar(&moodStr, "mood", "", "Filter by mood")
	cmd.Flags().StringVar(&tag, "tag", "", "Filter by tag")
	return cmd
}

func getMoodEmoji(m domain.Mood) string {
	switch m {
	case domain.MoodGreat:
		return "😊"
	case domain.MoodGood:
		return "🙂"
	case domain.MoodOkay:
		return "😐"
	case domain.MoodBad:
		return "🙁"
	case domain.MoodTerrible:
		return "😢"
	default:
		return "📝"
	}
}
