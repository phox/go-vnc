package vnc

import (
	"bytes"
	"net"
	"time"
)

// MockConn implements the net.Conn interface.
type MockConn struct {
	b bytes.Buffer
}

func (m *MockConn) Read(b []byte) (int, error) {
	return m.b.Read(b)
}
func (m *MockConn) Write(b []byte) (int, error) {
	return m.b.Write(b)
}
func (m *MockConn) Close() error                       { return nil }
func (m *MockConn) LocalAddr() net.Addr                { return nil }
func (m *MockConn) RemoteAddr() net.Addr               { return nil }
func (m *MockConn) SetDeadline(t time.Time) error      { return nil }
func (m *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *MockConn) SetWriteDeadline(t time.Time) error { return nil }

// Implement additional buffer.Buffer functions.
func (m *MockConn) Reset() {
	m.b.Reset()
}
