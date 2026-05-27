package commands

import (
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/logger"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewGetCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get <key>",
		Short:   "Get the value of a configuration key",
		Example: "Delete key of name from ./config.json:\ncfgctl --path ./config.json delete name",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandGet
			logger, err := logger.NewLogger(string(cmdKind))

			if err != nil {
				return err
			}

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
				commandError := NewCommandError(KindFileStoreLoad, cmdKind, err)
				logger.Error(commandError.DetailMessage())
				return commandError
			}

			key := args[0]
			value, err := fileStore.Get(key)

			if err != nil {
				commadError := NewCommandError(KindFileStoreGet, cmdKind, fmt.Errorf("get: failed to get value from key %s %w", key, err))
				logger.Error(commadError.Detail.Error())
				return commadError
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
