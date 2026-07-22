package cli

import (
	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/spf13/cobra"
)

func Setup(root *cobra.Command, cfg *config.Config) {
	fileRepo, _ := repository.NewFileRepository(cfg.StorageDir)

	//Add commands
	root.AddCommand(newCreateCmd(fileRepo))
	//root.AddCommand(newListCmd(fileRepo))
}
