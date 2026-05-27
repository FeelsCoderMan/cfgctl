package commands

import (
	"encoding/json"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewInitCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdKind := CommandInit
			path := cmd.Flag("path").Value.String()
			if path == "" {
				return NewCommandError(KindMissingPath, cmdKind, nil)
			}

			fileStore.SetPath(path)

			var jsonData map[string]any
			data, _ := cmd.Flags().GetString("data")

			if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
				return NewCommandError(KindInvalidJSON, cmdKind, err)
			}

			if err := fileStore.LoadData(jsonData); err != nil {
				return NewCommandError(KindFileStoreLoad, cmdKind, err)
			}

			if jsonData != nil {
				cmd.Printf("init: configuration is saved to %s with data\n", path)
			} else {
				cmd.Printf("init: configuration file is created to %s\n", path)
			}

			return nil
		},
	}

	cmd.Flags().String("data", "{}", "Initial configuration data")
	return cmd
}
