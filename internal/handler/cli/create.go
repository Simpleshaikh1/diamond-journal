package cli

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
)

func newCreateCmd(repo *repository.FileRepository) *cobra.Command {
	var moodStr, tagsStr string

	cmd := &cobra.Command{
		Use:   "new [title]",
		Short: "Create a new journal entry",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			title := ""
			if len(args) > 0 {
				title = strings.Join(args, " ")
			} else {
				fmt.Print("Title: ")
				title, _ = bufio.NewReader(os.Stdin).ReadString('\n')
				title = strings.TrimSpace(title)
			}

			fmt.Println("Write your entry (press Ctrl+D when done):")
			contentBytes, _ := os.ReadAll(os.Stdin)
			content := string(contentBytes)

			entry := &domain.Entry{
				Title:   title,
				Content: content,
				Mood:    domain.Mood(moodStr),
				Tags:    strings.Split(tagsStr, ","),
			}

			if err := repo.Create(entry); err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Printf("✅ Entry #%d created successfully!\n", entry.ID)

		},
	}

	cmd.Flags().StringVar(&moodStr, "mood", string(domain.MoodOkay), "mood")
	cmd.Flags().StringVar(&tagsStr, "tags", "", "comma-separated tags")
	return cmd
}
