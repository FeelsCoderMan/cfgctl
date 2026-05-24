package commands

import (
	"errors"
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewGetCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get the value of a configuration key",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, fmt.Errorf("get: missing path argument"))
			}

			if len(args) < 1 {
				return NewCommandError(KindMissingArgs, fmt.Errorf("get: missing key argument"))
			} else if len(args) > 1 {
				return NewCommandError(KindTooManyArgs, fmt.Errorf("get: too many arguments"))
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				if errors.Is(err, storage.ErrInvalidJSON) {
					return NewCommandError(KindInvalidJSON, err)
				}
				return NewCommandError(KindFileStoreLoad, fmt.Errorf("get: failed to load configuration file %s", path))
			}

			key := args[0]
			value, err := fileStore.Get(key)

			if err != nil {
				return NewCommandError(KindFileStoreGet, fmt.Errorf("get: failed to get value from key %s", key))
			}

			cmd.Printf("get: %s = %v\n", key, value)
			return nil
		},
	}

	return cmd
}
