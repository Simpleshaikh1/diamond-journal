package cli

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"

	"github.com/Simpleshaikh1/diamond-journal/internal/domain"
)

func newCreateCmd(repo domain.EntryRepository) *cobra.Command {
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
			content := readMultiLine()

			if content == "" {
				fmt.Println("⚠️  Empty content. Entry not saved.")
				return
			}

			entry := &domain.Entry{
				Title:     title,
				Content:   content,
				Mood:      domain.Mood(moodStr),
				Tags:      parseTags(tagsStr),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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

//func readMultiLine() string {
//	var lines []string
//	scanner := bufio.NewScanner(os.Stdin)
//	for scanner.Scan() {
//		line := scanner.Text()
//		if line == "" {
//			// Check if next line is also empty
//			if scanner.Scan() && scanner.Text() == "" {
//				break
//			}
//		}
//		lines = append(lines, line)
//	}
//	return strings.Join(lines, "\n")
//}

func readMultiLine() string {
	var content strings.Builder
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Warning: %v\n", err)
	}

	return strings.TrimSpace(content.String())
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}
	var tags []string
	for _, t := range strings.Split(tagsStr, ",") {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			tags = append(tags, trimmed)
		}
	}
	return tags
}
