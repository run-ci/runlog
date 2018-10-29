package cmd

import (
	"fmt"
	"os"

	"github.com/run-ci/runlog/pkg/runlog"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "GET /log/:taskID on runlogd",
	Long: `"get log" is for retrieving logs from runlogd.

It streams the logs back over the websocket connection until
the connection is closed.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := os.Getenv("RUNLOGD_URL")
		if addr == "" {
			addr = "localhost:7777"
		}

		user := os.Getenv("RUNLOGD_USER")
		if user == "" {
			fmt.Printf("need RUNLOGD_USER")
			os.Exit(1)
		}

		pass := os.Getenv("RUNLOGD_PASS")
		if pass == "" {
			fmt.Printf("need RUNLOGD_PASS")
			os.Exit(1)
		}

		client := runlog.NewClient(addr, user, pass)

		err := client.GetLog(1, os.Stdout)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(logCmd)
}
