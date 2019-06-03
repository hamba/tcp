package tcp_test

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/hamba/tcp"
)

type clientCodec struct {
	pool *ConnectionPool
	conn io.ReadWriter
}

func (c *clientCodec) Write(ctx context.Context, w io.Writer) error {
	panic("TODO")
}

func (c *clientCodec) Read(ctx context.Context, r io.Reader) error {
	panic("TODO")
}

func (c *clientCodec) Close() error {
	return nil
}

var DefaultPool = &ConnectionPool{
	IdleTimeout: 3 * time.Second,
}

type Pool interface {
	tcp.Pool
}

type ConnectionPool struct {
	tcp.Pool

	IdleTimeout time.Duration
}

type Client struct {
	client *tcp.Client
}

func NewClient(pool Pool, timeout time.Duration) (*Client, error) {
	if pool == nil {
		pool = DefaultPool
	}

	client := &tcp.Client{
		Timeout: timeout,
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Send(ctx context.Context, addr string, b []byte) ([]byte, error) {
	inBuf := bytes.NewBuffer(b)
	outBuf := bytes.NewBuffer(nil)
	if err := c.client.Send(ctx, nil, addr, inBuf, outBuf); err != nil {
		return nil, err
	}

	return outBuf.Bytes(), nil

}
