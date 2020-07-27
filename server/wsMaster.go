package server

import "github.com/gorilla/websocket"

type WsMaster struct {
	*websocket.Conn
}

func (w *WsMaster) Write(p []byte) (n int, err error) {
	writer, err := w.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	defer writer.Close()
	return writer.Write(p)
}

func (w *WsMaster) Read(p []byte) (n int, err error) {
	for {
		msgType, reader, err := w.Conn.NextReader()
		if err != nil {
			return 0, err
		}
		if msgType != websocket.TextMessage {
			continue
		}
		return reader.Read(p)
	}
}

func (w *WsMaster) Close() error {
	return w.Conn.Close()
}
