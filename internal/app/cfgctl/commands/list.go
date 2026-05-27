package commands

import (
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewListCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandList
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, cmdKind, nil)
			}

			if len(args) > 0 {
				return NewCommandError(KindTooManyArgs, cmdKind, nil)
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				return NewCommandError(KindFileStoreLoad, cmdKind, err)
			}

			result := fileStore.List()

			if len(result) == 0 {
				cmd.Println("list: no configuration values found")
				return nil
			} else {
				cmd.Println("list: configurations:")
			}

			for key, value := range result {
				cmd.Printf("-- %s: %v\n", key, value)
			}

			return nil
		},
	}

	return cmd
}
