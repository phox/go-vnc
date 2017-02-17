package vnc

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/kward/go-vnc/go/operators"
)

func TestBuffer_Read(t *testing.T) {
	var (
		buf    *Buffer
		tbyte  byte
		tint32 int32
	)

	buf = NewBuffer(nil)
	if err := buf.Read(&tint32); err == nil {
		t.Error("expected error")
	}

	buf = NewBuffer([]byte{123})
	if err := buf.Read(&tbyte); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := tbyte, byte(123); got != want {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}

	buf = NewBuffer([]byte{0, 18, 214, 135})
	if err := buf.Read(&tint32); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := tint32, int32(1234567); got != want {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}
}

func TestBuffer_Write(t *testing.T) {
	var buf *Buffer

	buf = NewBuffer(nil)
	if err := buf.Write(byte(234)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := buf.Bytes(), []byte{234}; !operators.EqualSlicesOfByte(got, want) {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}

	buf = NewBuffer(nil)
	if err := buf.Write(int32(23637)); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got, want := buf.Bytes(), []byte{0, 0, 92, 85}; !operators.EqualSlicesOfByte(got, want) {
		t.Errorf("incorrect result; got = %v, want = %v", got, want)
	}
}

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
