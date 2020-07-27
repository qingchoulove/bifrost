package tty

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"qingchoulove.github.com/bifrost/util"
	"sync"
)

type WebTTY struct {
	master     Master
	slave      Slave
	bufferSize int
	writeMutex sync.Mutex
}

func NewWebTTY(master Master, slave Slave) (*WebTTY, error) {
	tty := &WebTTY{
		master:     master,
		slave:      slave,
		bufferSize: 1024 * 64,
	}
	return tty, nil
}

func (tty *WebTTY) Run(ctx context.Context) error {
	quit := util.NewCloseChannel()
	defer quit.Close()
	// slave read
	go func() {
		buf := make([]byte, tty.bufferSize)
		for {
			n, err := tty.slave.Read(buf)
			if err != nil {
				log.Println("slave read:", err)
				quit.Close()
				return
			}

			n, err = tty.master.Write(encode(buf[:n]))
			if err != nil {
				log.Println("master write:", err)
				quit.Close()
				return
			}
		}
	}()
	// master read
	go func() {
		buf := make([]byte, tty.bufferSize)
		for {
			n, err := tty.master.Read(buf)
			if err != nil {
				log.Println("master read:", err)
				quit.Close()
				return
			}
			err = tty.handleMasterMessage(buf[:n])
			if err != nil {
				log.Println("slave write:", err)
				quit.Close()
				return
			}
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-quit.Wait():
		return errors.New("webTTY quit")
	}
}

func (tty *WebTTY) handleMasterMessage(payload []byte) error {
	if len(payload) == 0 {
		return errors.New("unexpected zero length read from master")
	}
	body := decode(payload[1:])
	switch payload[0] {
	case Input:
		_, err := tty.slave.Write(body)
		return err
	case ResizeTerminal:
		var args argResizeTerminal
		err := json.Unmarshal(body, &args)
		if err != nil {
			return err
		}
		return tty.slave.ResizeTerminal(args.Rows, args.Columns)
	default:
		return fmt.Errorf("unknown message type `%c`", payload[0])
	}
}

func decode(p []byte) []byte {
	decodeString, _ := base64.StdEncoding.DecodeString(string(p))
	return decodeString
}

func encode(p []byte) []byte {
	encodeToString := base64.StdEncoding.EncodeToString(p)
	return []byte(encodeToString)
}

type argResizeTerminal struct {
	Columns uint16
	Rows    uint16
}
