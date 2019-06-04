package tcp_test

import (
	"context"
	"io"

	"github.com/hamba/tcp"
)

type clientCodec struct {
	client *Client
	conn   io.ReadWriter
}

func (c *clientCodec) Write(ctx context.Context, w tcp.RequestWriter) error {
	// Need to write in protocol
	return w.Write(c.conn)
}

func (c *clientCodec) Read(ctx context.Context, r tcp.ResponseReader) error {
	// We got a response, read it and notify the client
	panic("TODO")
}

func (c *clientCodec) Close() error {
	return nil
}

type Client struct {
	client *tcp.Client
}

func NewClient(addr string) (*Client, error) {
	c := &Client{}

	client, err := tcp.NewClient(addr, c.createCodec, tcp.ClientOpts{})
	if err != nil {
		return nil, err
	}
	c.client = client

	return c, nil
}

func (c *Client) createCodec(conn tcp.Connection) tcp.ClientCodec {
	return &clientCodec{
		client: c,
		conn:   conn,
	}
}

func (c *Client) Call(ctx context.Context, addr string, r *Request) ([]byte, error) {
	select {
	case codec := <-c.client.Acquire():
		if err := codec.Write(ctx, r); err != nil {
			// Dont release the connection, close instead
			return nil, err
		}

	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Wait for read, or timeout

	return nil, nil
}
