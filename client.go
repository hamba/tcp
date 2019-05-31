package tcp

import (
	"context"
	"errors"
	"io"
	"net"
	"time"
)

type ClientCodec interface {
	Connection() Connection
	Write(ctx context.Context, w io.Writer) (read bool, err error)
	Read(ctx context.Context, r io.Reader) error
}

type ClientCodecFactory func(conn Connection) ClientCodec

type PoolOpts struct {
	// Dialer is the dialer to use
	Dialer net.Dialer

	// TLS Config

	// TLS Handshake Timeout

	// IdleTimeout is the maximum duration that a connection will remain idle for.
	IdleTimeout time.Duration
}

// Pool represents a connection pool
type Pool interface {
	Get(addr string) (Connection, ClientCodec, error)

	Put(Connection)

	unexported()
}

type pool struct {
	fac         ClientCodecFactory
	dialer      net.Dialer
	idleTimeout time.Duration
}

func NewPool(fac ClientCodecFactory, opts PoolOpts) (Pool, error) {
	if fac == nil {
		return nil, errors.New("tcp: fac cannot be nil")
	}

	return &pool{
		fac:         fac,
		dialer:      opts.Dialer,
		idleTimeout: opts.IdleTimeout,
	}, nil
}

func (p *pool) Get(addr string) (Connection, ClientCodec, error) {
	return nil, nil, nil
}

func (p *pool) Put(conn Connection) {

}

func (p *pool) unexported() {

}

type Client struct {
	// Timeout is the maximum duration to wait for a request to complete.
	Timeout time.Duration
}

func (c *Client) Send(ctx context.Context, pool Pool, addr string, w io.Writer, r io.Reader) error {
	if pool == nil {
		// Perhaps, dial and send end
		return errors.New("tcp: pool cannot be nil")
	}

	// All in a loop with checks for retry: conn could have been closed on idle, etc.

	// Get connection+codec from pool

	// Run the codec Write

	// If there is no response reader, we are done.

	// Wait for data, or timeout (tell the conn we are waiting for data)

	// Run the codec read

	for {
		conn, codec, err := pool.Get(addr)
		if err != nil {
			return err
		}


	}

	panic("TODO")
}
