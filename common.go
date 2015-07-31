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

func (b *Buffer) WriteByte(c byte) {
	b.buf.WriteByte(c)
	return
}

func (b *Buffer) WriteBytes(p []byte) {
	for _, c := range p {
		b.WriteByte(c)
	}
}

// type rfbMessage interface {
// 	Unarshaler interface
// }

// Marshaler is the interface for objects that can marshal themselves.
type Marshaler interface {
	Marshal() ([]byte, error)
}

// func Marshal(msg rfbMessage) ([]byte, error) {
// 	// Can the object marshal itself?
// 	if m, ok := msg.(Marshaler); ok {
// 		return m.Marshal()
// 	}
// 	return []byte{}, fmt.Errorf("Message %v unable to marshal itself.", msg)
// }

// Unmarshaler is the interface for objects that can marshal themselves.
type Unmarshaler interface {
	Unmarshal([]byte) error
}

// func Unmarshal(buf []byte, msg rfbMessage) error {
// 	// Can the object unmarshal itslef?
// 	if m, ok := msg.(Unmarshaler); ok {
// 		return m.Unmarshal(buf)
// 	}
// 	return fmt.Errorf("Message %v unable to unmarshal itself.", msg)
// }
