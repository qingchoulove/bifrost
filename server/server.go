package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"qingchoulove.github.com/bifrost/tty"
)

type Server struct {
	upgrader *websocket.Upgrader
	static   string
}

func NewServer(startPath string) (*Server, error) {
	return &Server{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}, static: startPath,
	}, nil
}

func (server *Server) Run(ctx context.Context) error {
	router := gin.New()
	router.LoadHTMLFiles(filepath.Join(server.static, "front/dist/index.html"))
	router.Static("/dist", path.Join(server.static, "front/dist"))

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/ws", server.handleWebsocket(ctx))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	errs := make(chan error)
	go func() {
		errs <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Println("Server Shutdown:", err)
		}
		return ctx.Err()
	case err := <-errs:
		return err
	}
}

func (server *Server) handleWebsocket(ctx context.Context) gin.HandlerFunc {

	return func(c *gin.Context) {
		conn, err := server.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(500, "websocket upgrade: %s", err)
			return
		}
		master := &WsMaster{conn}
		defer master.Close()
		slave, err := tty.CreateSlave(tty.Local)
		if err != nil {
			master.Write([]byte(err.Error()))
			return
		}
		defer slave.Close()
		webTTY, err := tty.NewWebTTY(master, slave)
		if err != nil {
			master.Write([]byte(err.Error()))
			return
		}
		err = webTTY.Run(ctx)
		log.Println("webTTY Stop:", err)
	}
}
