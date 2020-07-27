package tty

import "io"

type Slave interface {
	io.ReadWriteCloser

	ResizeTerminal(rows, cols uint16) error
}

type resizeFunction func(rows, cols uint16) error
