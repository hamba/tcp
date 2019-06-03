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
	Write(ctx context.Context, w io.Writer) error
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
	Get(ctx context.Context, addr string) (Connection, error)

	Put(Connection)

	getCodec(Connection) ClientCodec
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

func (p *pool) Get(ctx context.Context, addr string) (Connection, error) {
	return nil, nil
}

func (p *pool) Put(conn Connection) {

}

func (p *pool) getCodec(Connection) ClientCodec {
	return nil
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

	if w == nil {
		return errors.New("tcp: w cannot be nil")
	}

	// All in a loop with checks for retry: conn could have been closed on idle, etc.

	// Get connection+codec from pool

	// Run the codec Write

	// If there is no response reader, we are done.

	// Wait for data, or timeout (tell the conn we are waiting for data)

	// Run the codec read

	for {
		conn, err := pool.Get(ctx, addr)
		if err != nil {
			return err
		}
		codec := pool.getCodec(conn)

		if err = codec.Write(ctx, w); err != nil {
			_ = conn.Close()
			return err
		}

		if r == nil {
			pool.Put(conn)
			return nil
		}

		// TODO: how to read from conn in a channel,
		// Also, the connection really needs a read loop at this point
		// so it can detect idle timeout
		select {
		case err = <-ctx.Done():
			_ = conn.Close()
			return err
		}

		if err = codec.Read(ctx, r); err == nil {
			return nil
		}

		// TODO: should we try again
	}
}
