// Implementation of RFC 6143 §7.7 Encodings.

package vnc

import (
	"bytes"
	"fmt"
)

const (
	Raw               = int32(0)
	CopyRect          = int32(1)
	RRE               = int32(2)
	Hextile           = int32(5)
	TRLE              = int32(15)
	ZRLE              = int32(16)
	ColorPseudo       = int32(-239)
	DesktopSizePseudo = int32(-223)
)

// An Encoding implements a method for encoding pixel data that is
// sent by the server to the client.
type Encoding interface {
	// The number that uniquely identifies this encoding type.
	Type() int32

	// Read the contents of the encoded pixel data from the reader.
	// This should return a new Encoding implementation that contains
	// the proper data.
	Read(*ClientConn, *Rectangle) (Encoding, error)

	// Marshal implements the Marshaler interface.
	Marshal() ([]byte, error)
}

type Encodings []Encoding

func (e Encodings) Marshal() ([]byte, error) {
	buf := NewBuffer(nil)

	for _, enc := range e {
		if err := buf.Write(enc.Type()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// RawEncoding is raw pixel data sent by the server.
// See RFC 6143 §7.7.1.
type RawEncoding struct {
	Colors []Color
}

func (*RawEncoding) Type() int32 {
	return Raw
}

func (e *RawEncoding) Marshal() ([]byte, error) {
	buf := NewBuffer(nil)

	for _, c := range e.Colors {
		bytes, err := c.Marshal()
		if err != nil {
			return nil, err
		}
		buf.WriteBytes(bytes)
	}

	return buf.Bytes(), nil
}

func (*RawEncoding) Read(c *ClientConn, rect *Rectangle) (Encoding, error) {
	// if c.debug {
	// 	log.Println("RawEncoding.Read()")
	// 	rect.DebugPrint()
	// }

	var buf bytes.Buffer

	bytesPerPixel := int(c.pixelFormat.BPP / 8)
	n := rect.Area() * bytesPerPixel
	if err := c.receiveN(&buf, n); err != nil {
		return nil, fmt.Errorf("unable to read rectangle with raw encoding: %v", err)
	}

	colors := make([]Color, rect.Area())
	for y := uint16(0); y < rect.Height; y++ {
		for x := uint16(0); x < rect.Width; x++ {
			color := NewColor(&c.pixelFormat, &c.colorMap)
			if err := color.Unmarshal(buf.Next(bytesPerPixel)); err != nil {
				return nil, err
			}
			colors[int(y)*int(rect.Width)+int(x)] = *color
		}
	}

	return &RawEncoding{colors}, nil
}

// DesktopSizePseudoEncoding enables desktop resize support.
// See RFC 6143 §7.8.2.
type DesktopSizePseudoEncoding struct{}

func (*DesktopSizePseudoEncoding) Type() int32 {
	return DesktopSizePseudo
}

func (*DesktopSizePseudoEncoding) Read(c *ClientConn, rect *Rectangle) (Encoding, error) {
	c.fbWidth = rect.Width
	c.fbHeight = rect.Height

	return &DesktopSizePseudoEncoding{}, nil
}

func (e *DesktopSizePseudoEncoding) Marshal() ([]byte, error) {
	return []byte{}, nil
}
