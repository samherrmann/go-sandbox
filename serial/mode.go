package serial

import "go.bug.st/serial"

var mode = &serial.Mode{
	BaudRate: 38400,
	Parity:   serial.NoParity,
	DataBits: 8,
	StopBits: serial.OneStopBit,
}
