// Common things that aren't part of the RFB protocol.

package vnc

import (
	"bytes"
	"encoding/binary"
	"time"
)

const (
	RFBFalse = uint8(iota)
	RFBTrue
)

// VNCError implements error interface.
type VNCError struct {
	s string
}

// NewVNCError returns a custom VNCError error.
func NewVNCError(s string) error {
	return &VNCError{s}
}

// Error returns an VNCError as a string.
func (e VNCError) Error() string {
	return e.s
}

var settleDuration = 25 * time.Millisecond

// Settle returns the UI settle duration.
func Settle() time.Duration {
	return settleDuration
}

// SetSettle changes the UI settle duration.
func SetSettle(s time.Duration) {
	settleDuration = s
}

// settleUI allows the UI to "settle" before the next UI change is made.
func settleUI() {
	time.Sleep(settleDuration)
}

type Buffer struct {
	buf *bytes.Buffer // byte stream
}

func NewBuffer(b []byte) *Buffer {
	return &Buffer{buf: bytes.NewBuffer(b)}
}

func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *Buffer) Read(data interface{}) error {
	return binary.Read(b.buf, binary.BigEndian, data)
}

func (b *Buffer) Write(data interface{}) error {
	return binary.Write(b.buf, binary.BigEndian, data)
}

func (b *Buffer) WriteByte(c byte) error {
	return b.buf.WriteByte(c)
}

// Marshaler is the interface satisfied for un-/marshaling.
type Marshaler interface {
	// Marshal returns the wire encoding of the ServerMessage.
	Marshal() ([]byte, error)
	// Unmarshal parses a wire format message into a ServerMessage.
	Unmarshal(data []byte) error
}
