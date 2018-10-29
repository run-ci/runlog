package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/run-ci/runlog/pkg/runlog"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [TASK ID]",
	Short: "stream the log for the given task to stdout",
	Long: `"get log" is for retrieving logs from runlogd.

It streams the logs back over the websocket connection until
the connection is closed.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := os.Getenv("RUNLOGD_URL")
		if addr == "" {
			addr = "localhost:7777"
		}

		user := os.Getenv("RUNLOGD_USER")
		if user == "" {
			fmt.Println("need RUNLOGD_USER")
			os.Exit(1)
		}

		pass := os.Getenv("RUNLOGD_PASS")
		if pass == "" {
			fmt.Println("need RUNLOGD_PASS")
			os.Exit(1)
		}

		client := runlog.NewClient(addr, user, pass)

		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = client.GetLog(id, os.Stdout)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(logCmd)
}
