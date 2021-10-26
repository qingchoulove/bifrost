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
	"strconv"
)

type Server struct {
	ctx      context.Context
	upgrader *websocket.Upgrader
	static   string
	username string
	password string
	port     int
}

func NewServer(ctx context.Context, options ...Option) (*Server, error) {
	s := &Server{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ctx: ctx,
	}
	for _, option := range options {
		option(s)
	}
	return s, nil
}

func (server *Server) Run() error {
	router := gin.New()
	router.LoadHTMLFiles(filepath.Join(server.static, "front/dist/index.html"))
	router.Static("/dist", path.Join(server.static, "front/dist"))

	basicAuth := gin.BasicAuth(gin.Accounts{
		server.username: server.password,
	})
	router.GET("/", basicAuth, func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	router.GET("/ws", basicAuth, server.handleWebsocket(server.ctx))
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(server.port),
		Handler: router,
	}
	errs := make(chan error)
	go func() {
		errs <- srv.ListenAndServe()
	}()
	select {
	case <-server.ctx.Done():
		err := srv.Shutdown(server.ctx)
		if err != nil {
			log.Println("Server Shutdown:", err)
		}
		return server.ctx.Err()
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
