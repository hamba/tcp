package tcp_test

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"testing"
	"time"

	"github.com/hamba/tcp"
)

func newTestServer(t testing.TB, fac tcp.ServerCodecFactory, opts tcp.ServerOpts) (net.Addr, *tcp.Server) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	srv, err := tcp.NewServer(fac, opts)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err = srv.Serve(ln)
		if err != nil && err != tcp.ErrServerClosed {
			t.Fatal(err)
		}
	}()

	return ln.Addr(), srv
}

func TestNewServer_ErrorsOnNilFactory(t *testing.T) {
	_, err := tcp.NewServer(nil, tcp.ServerOpts{})

	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestServer_ServesConnectionCloses(t *testing.T) {
	addr, srv := newTestServer(t, func(conn tcp.Connection) tcp.ServerCodec {
		return &pingCodec{close: true, conn: conn}
	}, tcp.ServerOpts{})
	defer srv.Close()

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatal("dial error", err)
	}

	if _, err := io.WriteString(conn, "ping"); err != nil {
		t.Fatal("write error", err)
	}

	done := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(5 * time.Second):
			t.Error("body not closed after 5s")
			return
		case <-done:
		}
	}()

	if _, err := ioutil.ReadAll(conn); err != nil {
		t.Fatal("read error", err)
	}
	done <- true
}

func TestServer_ServesConnectionStaysOpen(t *testing.T) {
	addr, srv := newTestServer(t, func(conn tcp.Connection) tcp.ServerCodec {
		return &pingCodec{conn: conn}
	}, tcp.ServerOpts{})
	defer srv.Close()

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		t.Fatal("dial error", err)
	}
	defer conn.Close()

	for i := 0; i < 3; i++ {
		if _, err := io.WriteString(conn, "ping"); err != nil {
			t.Fatal("write error", err)
		}

		var pong [4]byte
		if _, err := conn.Read(pong[:]); err != nil || string(pong[:]) != "pong" {
			t.Fatal("read error", err)
		}
	}
}

type pingCodec struct {
	writeTimeout time.Duration
	close        bool

	conn tcp.Connection
}

func (c *pingCodec) Handle(ctx context.Context, deadline tcp.SetWriteDeadline) {
	var ping [4]byte
	c.conn.Read(ping[:])

	if c.writeTimeout > 0 {
		deadline(time.Now().Add(c.writeTimeout))
	}

	c.conn.Write([]byte("pong"))

	if c.close {
		c.conn.Close()
	}
}
