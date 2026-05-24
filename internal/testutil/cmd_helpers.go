package testutil

import (
	"bytes"

	"github.com/spf13/cobra"
)

func RunCmd(cmd *cobra.Command, args ...string) error {
	bufferOut := &bytes.Buffer{}
	bufferErr := &bytes.Buffer{}
	cmd.SetOut(bufferOut)
	cmd.SetErr(bufferErr)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func NewMockRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "root",
		Short: "main Cfgctl CLI root command",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	rootCmd.PersistentFlags().StringP("path", "p", "", "path of the configuration file")
	rootCmd.SilenceUsage = true
	rootCmd.MarkFlagRequired("path")
	return rootCmd
}
