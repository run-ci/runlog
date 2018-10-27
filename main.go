package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Println("opening file...")

	f, err := os.Open("logstash/tasklog/1.log")
	if err != nil {
		panic(err)
	}

	err = follow(os.Stdout, f)
	if err != nil {
		panic(err)
	}
}

func follow(w io.Writer, r io.Reader) error {
	buf := make([]byte, 32*1024)

	for {
		nr, err := r.Read(buf)
		if nr > 0 {
			if err != nil && err != io.EOF {
				return err
			}

			nw, err := w.Write(buf[0:nr])
			if err != nil {
				return err
			}

			if nr != nw {
				return io.ErrShortWrite
			}
		}
	}
}
