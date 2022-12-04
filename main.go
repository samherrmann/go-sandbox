package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/samherrmann/go-sandbox/serial"
)

const (
	ModeReader = "reader"
	ModeWriter = "writer"
)

type Config struct {
	Mode     string
	PortName string
}

func main() {
	err := app()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}

func app() error {
	// Parse arguments:
	c, err := parseArgs()
	if err != nil {
		programName := os.Args[0]
		usage :=
			fmt.Sprintf("Usage: %s <mode> <port-name>\n", programName) +
				"  mode       possible values: reader, writer\n" +
				"  port-name  name of the port\n\n"
		return fmt.Errorf("%s\n\n%s", err, usage)
	}

	// Catch interrupt signal to cleanly close reader/writer.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	// Create reader or writer:
	switch c.Mode {
	case ModeReader:
		return serial.Reader(c.PortName, done)
	case ModeWriter:
		return serial.Writer(c.PortName, done)
	default:
		return fmt.Errorf("unknown mode '%s'", c.Mode)
	}
}

func parseArgs() (*Config, error) {
	// Arguments without program name.
	args := os.Args[1:]
	if len(args) < 2 {
		return nil, errors.New("not enough arguments")
	}
	if len(args) > 2 {
		return nil, fmt.Errorf("unknown arguments: %s", os.Args[3:])
	}
	c := &Config{
		Mode:     args[0],
		PortName: args[1],
	}
	return c, nil
}
