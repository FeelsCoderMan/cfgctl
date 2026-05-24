package commands

import (
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a configuration value by key",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, fmt.Errorf("delete: missing path argument"))
			}

			if len(args) < 1 {
				return NewCommandError(KindMissingArgs, fmt.Errorf("delete: missing key argument"))
			} else if len(args) > 1 {
				return NewCommandError(KindTooManyArgs, fmt.Errorf("delete: too many arguments"))
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				if errors.Is(err, storage.ErrInvalidJSON) {
					return NewCommandError(KindInvalidJSON, err)
				}
				return NewCommandError(KindFileStoreLoad, fmt.Errorf("delete: failed to load configuration file %s", path))
			}

			key := args[0]

			if err := fileStore.Delete(key); err != nil {
				return NewCommandError(KindFileStoreDelete, fmt.Errorf("delete: failed to delete key %s", key))
			}

			cmd.Printf("delete: key %s is deleted from configuration file\n", key)
			return nil
		},
	}

	return cmd
}
