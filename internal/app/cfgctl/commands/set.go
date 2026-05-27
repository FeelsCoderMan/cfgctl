package commands

import (
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewSetCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandSet
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, cmdKind, nil)
			}

			if len(args) < 2 {
				return NewCommandError(KindMissingArgs, cmdKind, nil)
			} else if len(args) > 2 {
				return NewCommandError(KindTooManyArgs, cmdKind, nil)
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				return NewCommandError(KindFileStoreLoad, cmdKind, fmt.Errorf("set: failed to load configuration file %s: %w", path, err))
			}

			key, value := args[0], args[1]

			if err := fileStore.Set(key, value); err != nil {
				return NewCommandError(KindFileStoreSet, cmdKind, fmt.Errorf("set: failed to set key %s to value %s: %w", key, value, err))
			}

			cmd.Printf("set: key %s set to value %s\n", key, value)
			return nil
		},
	}

	return cmd
}
