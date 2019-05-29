package tcp_test

import (
	"context"
	"io"

	"github.com/hamba/tcp"
)

type clientCodec struct {
	client *Client

	conn io.ReadWriter
}

func (c *clientCodec) Write(ctx context.Context, w tcp.RequestWriter) (read bool, err error) {
	panic("TODO")
}

func (c *clientCodec) Read(ctx context.Context, r tcp.ResponseReader) error {
	panic("TODO")
}

func (c *clientCodec) Close() error {
	return nil
}

type Client struct {
	client *tcp.Client
}

func (c *Client) createCodec(conn tcp.Connection) tcp.ClientCodec {
	return &clientCodec{
		client: c,
		conn:   conn,
	}
}

func (c *Client) Send(ctx context.Context, ) {

}
