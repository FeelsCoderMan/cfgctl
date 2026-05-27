package commands

import (
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewGetCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the value of a configuration key",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandGet
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
				return NewCommandError(KindFileStoreLoad, cmdKind, err)
			}

			key := args[0]
			value, err := fileStore.Get(key)

			if err != nil {
				return NewCommandError(KindFileStoreGet, cmdKind, fmt.Errorf("get: failed to get value from key %s %w", key, err))
			}

			if value == nil {
				cmd.Printf("get: key %s not found\n", key)
				return nil
			}

			cmd.Printf("get: %s = %v\n", key, value)
			return nil
		},
	}

	return cmd
}
