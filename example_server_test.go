package tcp_test

import (
	"context"
	"errors"
	"io"
	"net"
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

type serverCodec struct {
	server *Server

	conn io.ReadWriter
}

func (c *serverCodec) Handle(ctx context.Context, deadline tcp.SetWriteDeadline) bool {
	req := &Request{Body: c.conn}

	deadline(time.Now().Add(c.server.WriteTimeout))

	c.server.Handler.ServeTCP(c.conn, req)

	return req.Close
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

	return s.srv.Serve(ln)
}

func (s *Server) createCodec(conn io.ReadWriter) tcp.ServerCodec {
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
