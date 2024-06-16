package mock

import "io"

type Conn struct {
	OnRead  func([]byte) (int, error)
	OnWrite func([]byte) (int, error)
	OnClose func() error
}

var _ io.ReadWriteCloser = (*Conn)(nil)

func (c *Conn) Read(p []byte) (n int, err error) {
	if c.OnRead == nil {
		return 0, nil
	}
	return c.OnRead(p)
}

func (c *Conn) Write(p []byte) (n int, err error) {
	if c.OnWrite == nil {
		return 0, nil
	}
	return c.OnWrite(p)
}

func (c *Conn) Close() error {
	if c.OnClose == nil {
		return nil
	}
	return c.OnClose()
}
