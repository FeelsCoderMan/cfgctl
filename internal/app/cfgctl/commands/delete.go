package commands

import (
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a configuration value by key",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandDelete
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, cmdKind, nil)
			}

			if len(args) < 1 {
				return NewCommandError(KindMissingArgs, cmdKind, nil)
			} else if len(args) > 1 {
				return NewCommandError(KindTooManyArgs, cmdKind, nil)
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				commandError := NewCommandError(KindFileStoreLoad, cmdKind, fmt.Errorf("delete: failed to load configuration file %s: %w", path, err))
				return commandError
			}

			key := args[0]

			if err := fileStore.Delete(key); err != nil {
				commandError := NewCommandError(KindFileStoreDelete, cmdKind, fmt.Errorf("delete: failed to delete key %s: %w", key, err))
				return commandError
			}

			cmd.Printf("delete: key %s is deleted from configuration file\n", key)
			return nil
		},
	}

	return cmd
}
