package tcp_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/hamba/tcp"
)

type ResponseWriter interface {
	io.Writer
}

type Request struct {
	Body io.Reader

	Close bool
}

type Handler interface {
	ServeTCP(w ResponseWriter, r *Request)
}

type HandlerFunc func(w ResponseWriter, r *Request)

func (fn HandlerFunc) ServeTCP(w ResponseWriter, r *Request) {
	fn(w, r)
}

type serverCodec struct {
	server *Server

	conn io.ReadWriteCloser
}

func (c *serverCodec) Handle(ctx context.Context, deadline tcp.SetWriteDeadline) {
	req := &Request{Body: c.conn}

	if d := c.server.WriteTimeout; d > 0 {
		deadline(time.Now().Add(d))
	}

	c.server.Handler.ServeTCP(c.conn, req)

	if req.Close {
		_ = c.conn.Close()
	}
}

func (c *serverCodec) Close() error {
	return nil
}

type Server struct {
	Addr string

	Handler Handler

	ReadTimeout time.Duration

	WriteTimeout time.Duration

	IdleTimeout time.Duration

	mu  sync.Mutex
	srv *tcp.Server
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":8080"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(ln)
}

func (s *Server) Serve(ln net.Listener) error {
	if s.Handler == nil {
		return errors.New("server: handler cannot be nil")
	}

	s.mu.Lock()
	if s.srv == nil {
		srv, err := tcp.NewServer(s.createCodec, tcp.ServerOpts{
			ReadTimeout: s.ReadTimeout,
			IdleTimeout: s.IdleTimeout,
		})
		if err != nil {
			return err
		}

		s.srv = srv
	}
	s.mu.Unlock()

	return s.srv.Serve(ln)
}

func (s *Server) createCodec(conn tcp.Connection) tcp.ServerCodec {
	return &serverCodec{
		server: s,
		conn:   conn,
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Close() error {
	return s.srv.Close()
}

func ExampleServer() {
	h := HandlerFunc(func(w ResponseWriter, r *Request) {
		// Do your read and write here
	})

	srv := &Server{
		Addr:         "localhost:8090",
		Handler:      h,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
		IdleTimeout:  10 * time.Second,
	}
	defer srv.Close()

	if err := srv.ListenAndServe(); err != nil {
		log.Print(err)
	}
}
