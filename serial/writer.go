package serial

import (
	"fmt"
	"os"
	"time"

	"go.bug.st/serial"
)

func Writer(portName string, done <-chan os.Signal) error {
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
			msg := time.Now().String()
			fmt.Printf("Sending: %s\n", msg)
			_, err = port.Write([]byte(msg))
			if err != nil {
				return err
			}
			time.Sleep(time.Second)
		}
	}
}
