package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "GET resources from runlogd",
	Long: `
"get" retrieves resources from runlogd. It only works with subcommands.
	`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
