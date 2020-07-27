package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"qingchoulove.github.com/bifrost/server"
	"syscall"
)

func main() {
	cancel, cancelFunc := context.WithCancel(context.Background())
	// graceful shutdown
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		log.Println("Caught signal:", s)
		cancelFunc()
	}()
	// run server
	serve, err := server.NewServer(getCurrPath())
	if err != nil {
		log.Println("NewServer Server:", err)
		return
	}
	err = serve.Run(cancel)
	if err != nil {
		log.Println("Server Stop:", err)
	}
}

func getCurrPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return path.Dir("")
	}
	return dir
}
