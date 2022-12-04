package serial

import (
	"fmt"
	"os"

	"go.bug.st/serial"
)

func Reader(portName string, done <-chan os.Signal) error {
	port, err := serial.Open(portName, mode)
	if err != nil {
		return err
	}
	defer port.Close()

	for {
		select {
		case <-done:
			return nil
		default:
			buf := make([]byte, 1024)
			_, err = port.Read(buf)
			if err != nil {
				return err
			}
			fmt.Printf("Received: %s\n", buf)
		}
	}
}
