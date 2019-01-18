package main

import (
	"fmt"
	"io"
	"os"

	"github.com/run-ci/runlog"
	flag "github.com/spf13/pflag"
)

var client *runlog.Client

func init() {
	url := os.Getenv("RUNLOGGER_TARGET_URL")
	if url == "" {
		fmt.Println("RUNLOGGER_TARGET not set, using localhost:9999")
		url = "localhost:9999"
	}

	capath := os.Getenv("RUNLOGGER_CAPATH")
	if capath == "" {
		fmt.Println("RUNLOGGER_CAPATH not set, using ./ca.pem")
		capath = "./ca.pem"
	}

	certpath := os.Getenv("RUNLOGGER_CERTPATH")
	if certpath == "" {
		fmt.Println("RUNLOGGER_CERTPATH not set, using ./runlog.crt")
		certpath = "./runlog.crt"
	}

	keypath := os.Getenv("RUNLOGGER_KEYPATH")
	if keypath == "" {
		fmt.Println("RUNLOGGER_KEYPATH not set, using ./runlog.key")
		keypath = "./runlog.key"
	}

	client = &runlog.Client{
		URL:      url,
		CertPath: certpath,
		KeyPath:  keypath,
		CAPath:   capath,
	}
}

func main() {
	fmt.Println("connecting to runlog server...")

	err := client.Connect()
	if err != nil {
		fmt.Printf("got error connecting: %v\n", err)
		os.Exit(1)
	}

	var task int32
	flag.Int32VarP(&task, "task", "t", 0, "the task id to send logs for")

	flag.Parse()

	client.TaskID = uint32(task)

	n, err := io.Copy(client, os.Stdin)
	if err != nil {
		fmt.Printf("got error copying stdin to client: %v\n", err)
		os.Exit(1)
	}

	//n, err := client.Write([]byte("Hello, World!"))
	//if err != nil {
	//fmt.Printf("got error writing message: %v\n", err)
	//os.Exit(1)
	//}

	fmt.Printf("wrote %v bytes\n", n)
}
