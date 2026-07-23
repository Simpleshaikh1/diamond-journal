package cli

import (
	"github.com/Simpleshaikh1/diamond-journal/internal/config"
	"github.com/Simpleshaikh1/diamond-journal/internal/repository"
	"github.com/spf13/cobra"
)

func Setup(root *cobra.Command, cfg *config.Config) {
	fileRepo, err := repository.NewFileRepository(cfg.StorageDir)
	if err != nil {
		panic(err)
	}

	//Add commands
	root.AddCommand(newCreateCmd(fileRepo))
	root.AddCommand(newListCmd(fileRepo))
	root.AddCommand(newViewCmd(fileRepo))
	root.AddCommand(newEditCmd(fileRepo))
	root.AddCommand(newDeleteCmd(fileRepo))
	root.AddCommand(newSearchCmd(fileRepo))
}
