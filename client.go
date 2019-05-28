package tcp

import (
	"context"
	"errors"
	"io"
	"time"
)

type RequestWriter interface {
	Write(io.Writer) error
}

type ResponseReader interface {
	Read(r io.Reader) error
}

type ClientCodec interface {
	Write(ctx context.Context, w RequestWriter) (read bool, err error)
	Read(ctx context.Context, r ResponseReader) error
}

type ClientCodecFactory func(conn Connection) ClientCodec

var DefaultPool Pool

type Pool interface {
	Get(addr string) (Connection, error)
}

type ClientOpts struct {
	// Timeout is the maximum duration to wait for a request to complete.
	Timeout time.Duration

	// Pool is the connection pool to use to get a connection.
	Pool Pool
}

type Client struct {
	fac  ClientCodecFactory
	pool Pool
}

func NewClient(fac ClientCodecFactory, opts ClientOpts) (*Client, error) {
	if fac == nil {
		return nil, errors.New("tcp: fac cannot be nil")
	}

	pool := DefaultPool
	if opts.Pool == nil {
		pool = opts.Pool
	}

	return &Client{
		fac:  fac,
		pool: pool,
	}, nil
}

func (c *Client) Send(ctx context.Context, w RequestWriter, r ResponseReader) error {
	// All in a loop with checks for retry: conn could have been closed on idle, etc.

	// Get connection+codec from pool

	// Run the codec Write

	// If there is no response reader, we are done.

	// Wait for data, or timeout (tell the conn we are waiting for data)

	// Run the codec read

	panic("TODO")
}
