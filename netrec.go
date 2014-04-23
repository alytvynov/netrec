// Package netrec provides a wrapper for net.Listener that
// records all data sent on it's connections
//
// Warning: all data is stored in memory while connection is alive,
// do not use this for large volumes of data!
package netrec

import (
	"bytes"
	"net"
)

// Callback is a function that gets triggered when connection is
// closed and receives all data that passed trough that connection.
// in contains everything read from connection;
// out contains everything written to connection.
type Callback func(in *bytes.Buffer, out *bytes.Buffer)

// NewRecordListener returns a new net.Listener, wrapping the one passed.
// Each call of Close on connections triggers cb. cb can not be nil.
func NewRecordListener(l net.Listener, cb Callback) net.Listener {
	return &recordListener{Listener: l, cb: cb}
}

type recordListener struct {
	net.Listener
	cb Callback
}

func (rl *recordListener) Accept() (net.Conn, error) {
	c, err := rl.Listener.Accept()
	return &recordConn{Conn: c, in: &bytes.Buffer{}, out: &bytes.Buffer{}, cb: rl.cb}, err
}

type recordConn struct {
	net.Conn
	in  *bytes.Buffer
	out *bytes.Buffer
	cb  Callback
}

func (rc *recordConn) Read(buf []byte) (int, error) {
	n, err := rc.Conn.Read(buf)
	// bytes.Buffer can't return an error, ignore
	rc.in.Write(buf[:n])
	return n, err
}

func (rc *recordConn) Write(buf []byte) (int, error) {
	// bytes.Buffer can't return an error, ignore
	rc.out.Write(buf)
	return rc.Conn.Write(buf)
}

func (rc *recordConn) Close() error {
	rc.cb(rc.in, rc.out)
	return rc.Conn.Close()
}
