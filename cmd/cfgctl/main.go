package main

import (
	"github.com/FeelsCoderMan/cfgctl/internal/app/cfgctl/commands"
	"github.com/FeelsCoderMan/cfgctl/internal/pkg/storage"
	"github.com/spf13/cobra"
)

func emptyRunFunc(cmd *cobra.Command, args []string) {}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cfgctl",
		Short: "cfgctl is a CLI tool for managing configuration files",
		Run:   emptyRunFunc,
	}

	rootCmd.PersistentFlags().StringP("path", "p", "", "path of the configuration file")
	rootCmd.MarkPersistentFlagRequired("path")

	fileStorage := storage.NewFileStore()

	rootCmd.AddCommand(
		commands.NewInitCmd(fileStorage),
		commands.NewGetCmd(fileStorage),
		commands.NewSetCmd(fileStorage),
		commands.NewDeleteCmd(fileStorage),
		commands.NewListCmd(fileStorage),
	)

	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln("Root command execution is failed ", err)
	}
}
