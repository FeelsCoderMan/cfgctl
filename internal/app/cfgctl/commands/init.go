package commands

import (
	"encoding/json"
	"fmt"

	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func NewInitCmd(fileStore *storage.FileStore) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			path := cmd.Flag("path").Value.String()
			if path == "" {
				return NewCommandError(KindMissingPath, fmt.Errorf("init: missing required flag ---path"))
			}

			fileStore.SetPath(path)

			var jsonData map[string]any
			data, _ := cmd.Flags().GetString("data")

			if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
				return NewCommandError(KindInvalidJSON, fmt.Errorf("init: invalid JSON in --data"))
			}

			if err := fileStore.LoadData(jsonData); err != nil {
				return NewCommandError(KindFileStoreLoad, fmt.Errorf("init: failed to load configuration from %s", path))
			}

			cmd.Println("init: configuration initialized successfully")
			return nil
		},
	}

	cmd.Flags().String("data", "{}", "Initial configuration data")
	return cmd
}
