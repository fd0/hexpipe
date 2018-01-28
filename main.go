package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

func hexdump(c *color.Color, in io.Reader, out io.WriteCloser, err chan<- error) {
	buf := make([]byte, 32*1024)
	for {
		n, err := in.Read(buf)
		data := buf[:n]

		if len(data) > 0 {
			s := hex.Dump(data)
			c.Fprint(os.Stderr, s)
		}

		_, werr := out.Write(data)
		if werr != nil {
			fmt.Fprintf(os.Stderr, "error writing: %v", werr)
			break
		}

		if err == io.EOF {
			c.Fprintln(os.Stderr, "EOF")
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading: %v", werr)
			break
		}
	}

	err <- out.Close()
}

func main() {
	color.NoColor = false

	command := os.Args[1]
	args := os.Args[2:]

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	errCh := make(chan error, 2)

	go hexdump(color.New(color.FgGreen), os.Stdin, stdin, errCh)
	go hexdump(color.New(color.FgRed), stdout, os.Stdout, errCh)

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(3)
	}
}
