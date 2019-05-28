package tcp

import "io"

// Connection represents a network connection.
type Connection interface {
	io.ReadWriteCloser

	// RemoteAddr returns the remote address.
	RemoteAddr() string

	//TLSInfo
}
