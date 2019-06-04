package tcp

import (
	"context"
	"errors"
	"io"
	"net"
)

type RequestWriter interface {
	Write(io.Writer) error
}

type ResponseReader interface {
	Read(r io.Reader) error
}

type ClientCodec interface {
	Write(ctx context.Context, w RequestWriter) error
	Read(ctx context.Context, r ResponseReader) error
}

type ClientCodecFactory func(conn Connection) ClientCodec

type ClientOpts struct {
	Dialer net.Dialer
}

type Client struct {
	connCh chan ClientCodec
}

func NewClient(addr string, fac ClientCodecFactory, opts ClientOpts) (*Client, error) {
	if fac == nil {
		return nil, errors.New("tcp: fac cannot be nil")
	}

	return &Client{
		connCh: make(chan ClientCodec, 1),
	}, nil
}

func (c *Client) Acquire() <-chan ClientCodec {
	return nil
}

func (c *Client) Close() error {
	// Bleed all waiting things

	return nil
}
