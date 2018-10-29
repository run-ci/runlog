package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "runlogq",
	Short: "runlogq is a client for runlog",
	Long:  `TODO: fill me out`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("hello client")
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
