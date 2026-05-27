package commands

import (
	"github.com/FeelsCoderMan/cfgctl/internal/logger"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewListCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all configuration values",
		Example: "List all configuration values from ./config.json:\ncfgctl list --path ./config.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandList
			logger, err := logger.NewLogger(string(cmdKind))

			if err != nil {
				return err
			}

			path := cmd.Flag("path").Value.String()

			if path == "" {
				return NewCommandError(KindMissingPath, cmdKind, nil)
			}

			if len(args) > 0 {
				return NewCommandError(KindTooManyArgs, cmdKind, nil)
			}

			fileStore.SetPath(path)

			if err := fileStore.LoadFromPath(); err != nil {
				commandError := NewCommandError(KindFileStoreLoad, cmdKind, err)
				logger.Error(commandError.DetailMessage())
				return commandError
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
