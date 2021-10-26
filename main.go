package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"qingchoulove.github.com/bifrost/server"
	"syscall"
)

var (
	cliPwd  = flag.String("pwd", "password", "Input Basic Auth Password")
	cliUser = flag.String("user", "user", "Input Basic Auth Username")
	cliPort = flag.Int("port", 8080, "Input web server port")
)

func main() {
	flag.Parse()
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
	serve, err := server.NewServer(cancel,
		server.OptionStaticPath(getCurrPath()),
		server.OptionUserName(*cliUser),
		server.OptionPassword(*cliPwd),
		server.OptionPort(*cliPort))
	if err != nil {
		log.Println("NewServer Server:", err)
		return
	}
	err = serve.Run()
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
