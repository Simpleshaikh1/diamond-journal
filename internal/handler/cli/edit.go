package cli

import (
	"fmt"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func newEditCmd(repo *repository.FileRepository) *cobra.Command {
	var useEditor bool

	cmd := &cobra.Command{
		Use:   "edit [id]",
		Short: "Edit an existing entry",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("❌ Invalid ID")
				return
			}

			entry, err := repo.GetByID(uint(id))
			if err != nil {
				fmt.Printf("❌ Entry #%d not found\n", id)
				return
			}

			fmt.Printf("Editing Entry #%d: %s\n", entry.ID, entry.Title)

			// Create temp file with current content
			tmpFile, _ := os.CreateTemp("", "journal-edit-*.md")
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					fmt.Printf("Failed to create temp file")
				}
			}(tmpFile.Name())

			os.WriteFile(tmpFile.Name(), []byte(entry.Content), 0644)

			if useEditor {
				editorCmd := exec.Command("code", "--wait", tmpFile.Name())
				editorCmd.Run()
			} else {
				fmt.Println("Enter new content (double Enter to finish):")
				newContent := readMultiLine()
				entry.Content = newContent
			}

			// Read updated content if editor was used
			if useEditor {
				newContent, _ := os.ReadFile(tmpFile.Name())
				entry.Content = string(newContent)
			}

			if strings.TrimSpace(entry.Content) == "" {
				fmt.Println("⚠️ Content cannot be empty.")
				return
			}

			entry.UpdatedAt = time.Now()

			if err := repo.Update(entry); err != nil {
				fmt.Printf("❌ Failed to update: %v\n", err)
				return
			}

			fmt.Printf("✅ Entry #%d updated successfully!\n", entry.ID)
		},
	}

	cmd.Flags().BoolVarP(&useEditor, "edit", "e", true, "Use external editor")
	return cmd
}
