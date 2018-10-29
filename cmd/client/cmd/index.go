package cmd

import (
	"fmt"
	"os"

	"github.com/run-ci/runlog/pkg/runlog"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "GET / on runlogd",
	Long: `
"get index" is for hitting the / route with GET.

This doesn't really return anything useful. It's only intended use
is to see if the server is even running or accessible.
`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := os.Getenv("RUNLOGD_URL")
		if addr == "" {
			addr = "http://localhost:7777"
		}

		// This request doesn't require auth.
		client := runlog.NewClient(addr, "", "")

		err := client.GetRoot()
		if err != nil {
			fmt.Printf("got error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(indexCmd)
}
