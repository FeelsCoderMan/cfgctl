package commands

import (
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewSetCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set a configuration value",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, fmt.Errorf("set: missing path"))
			}

			if len(args) < 2 {
				return NewCommandError(KindMissingArgs, fmt.Errorf("set: missing key and/or value"))
			} else if len(args) > 2 {
				return NewCommandError(KindTooManyArgs, fmt.Errorf("set: too many arguments"))
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				if errors.Is(err, storage.ErrInvalidJSON) {
					return NewCommandError(KindInvalidJSON, err)
				}
				return NewCommandError(KindFileStoreLoad, fmt.Errorf("set: failed to load configuration file %s", path))
			}

			key, value := args[0], args[1]

			if err := fileStore.Set(key, value); err != nil {
				return NewCommandError(KindFileStoreSet, fmt.Errorf("set: failed to set key %s to value %s", key, value))
			}

			cmd.Printf("set: key %s set to value %s\n", key, value)
			return nil
		},
	}

	return cmd
}
