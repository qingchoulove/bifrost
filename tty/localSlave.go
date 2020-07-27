package tty

import "os"

type LocalSlave struct {
	*os.File
	ResizeFunction resizeFunction
}

func (slave *LocalSlave) Read(p []byte) (n int, err error) {
	return slave.File.Read(p)
}

func (slave *LocalSlave) Write(p []byte) (n int, err error) {
	return slave.File.Write(p)
}

func (slave *LocalSlave) Close() error {
	return slave.File.Close()
}

func (slave *LocalSlave) ResizeTerminal(rows, cols uint16) error {
	return slave.ResizeFunction(rows, cols)
}
