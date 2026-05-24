package commands

import (
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewListCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configuration values",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, fmt.Errorf("list: path is required"))
			}

			if len(args) > 0 {
				return NewCommandError(KindTooManyArgs, fmt.Errorf("list: too many arguments"))
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				if errors.Is(err, storage.ErrInvalidJSON) {
					return NewCommandError(KindInvalidJSON, err)
				}
				return NewCommandError(KindFileStoreLoad, fmt.Errorf("list: failed to load configuration file %s", path))
			}

			result := fileStore.List()

			if len(result) == 0 {
				cmd.Println("list: No configuration values found")
				return nil
			} else {
				cmd.Println("list: Configurations:")
			}

			for key, value := range result {
				cmd.Printf("\t%s: %v\n", key, value)
			}

			return nil
		},
	}

	return cmd
}
