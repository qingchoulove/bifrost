package tty

import (
	"errors"
	"github.com/creack/pty"
	"os"
	"os/exec"
)

type SlaveType int

const (
	Local = iota
	Docker
)

func CreateSlave(slaveType SlaveType) (Slave, error) {
	switch slaveType {
	case Local:
		return createLocalSlave()
	default:
		return nil, errors.New("not support")
	}
}

func createLocalSlave() (*LocalSlave, error) {
	shell := os.Getenv("SHELL")
	if len(shell) == 0 {
		shell = "sh"
	}
	cmd := exec.Command(shell)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}
	return &LocalSlave{
		File: ptmx,
		ResizeFunction: func(rows, cols uint16) error {
			return pty.Setsize(ptmx, &pty.Winsize{
				Rows: rows,
				Cols: cols,
			})
		},
	}, nil
}
