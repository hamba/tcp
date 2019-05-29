package tcp

import (
	"context"
	"errors"
	"io"
	"net"
	"time"
)

type ClientCodec interface {
	Write(ctx context.Context, w io.Writer) (read bool, err error)
	Read(ctx context.Context, r io.Reader) error
}

type ClientCodecFactory func(conn Connection) ClientCodec

type PoolOpts struct {
	// Dialer is the dialer to use
	Dialer net.Dialer

	// IdleTimeout is the maximum duration that a connection will remain idle for.
	IdleTimeout time.Duration
}

// Pool represents a connection pool
type Pool interface {
	Get(addr string) (Connection, error)

	Put(Connection)
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

func (p *pool) Get(addr string) (Connection, error) {
	return nil, nil
}

func (p *pool) Put(conn Connection) {

}

type ClientOpts struct {
	// Timeout is the maximum duration to wait for a request to complete.
	Timeout time.Duration
}

type Client struct {
	pool Pool
}

func NewClient(pool Pool, opts ClientOpts) (*Client, error) {
	if pool == nil {
		return nil, errors.New("tcp: pool cannot be nil")
	}

	return &Client{

		pool: pool,
	}, nil
}

func (c *Client) Send(ctx context.Context, w io.Writer, r io.Reader) error {
	// All in a loop with checks for retry: conn could have been closed on idle, etc.

	// Get connection+codec from pool

	// Run the codec Write

	// If there is no response reader, we are done.

	// Wait for data, or timeout (tell the conn we are waiting for data)

	// Run the codec read

	panic("TODO")
}
