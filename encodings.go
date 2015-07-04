/*
encodings.go implements RFC 6143 §7.7 Encodings.
See http://tools.ietf.org/html/rfc6143#section-7.7 for more info.
*/
package vnc

import (
	"encoding/binary"
	"io"
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
	Read(*ClientConn, *Rectangle, io.Reader) (Encoding, error)
}

// RawEncoding is raw pixel data sent by the server.
// See RFC 6143 Section 7.7.1
type RawEncoding struct {
	Colors []Color
}

// NewRawEncoding returns a new RawEncoding.
func NewRawEncoding(c []Color) *RawEncoding {
	return &RawEncoding{c}
}

func (e *RawEncoding) Type() int32 {
	return Raw
}

func (*RawEncoding) Read(c *ClientConn, rect *Rectangle, r io.Reader) (Encoding, error) {
	bytesPerPixel := c.pixelFormat.BPP / 8
	pixelBytes := make([]uint8, bytesPerPixel)

	var byteOrder binary.ByteOrder = binary.LittleEndian
	if c.pixelFormat.BigEndian == RFBTrue {
		byteOrder = binary.BigEndian
	}

	colors := make([]Color, int(rect.Height)*int(rect.Width))
	for y := uint16(0); y < rect.Height; y++ {
		for x := uint16(0); x < rect.Width; x++ {
			if _, err := io.ReadFull(r, pixelBytes); err != nil {
				return nil, err
			}

			var rawPixel uint32
			if c.pixelFormat.BPP == 8 {
				rawPixel = uint32(pixelBytes[0])
			} else if c.pixelFormat.BPP == 16 {
				rawPixel = uint32(byteOrder.Uint16(pixelBytes))
			} else if c.pixelFormat.BPP == 32 {
				rawPixel = byteOrder.Uint32(pixelBytes)
			}

			color := &colors[int(y)*int(rect.Width)+int(x)]
			if c.pixelFormat.TrueColor == RFBTrue {
				color.R = uint16((rawPixel >> c.pixelFormat.RedShift) & uint32(c.pixelFormat.RedMax))
				color.G = uint16((rawPixel >> c.pixelFormat.GreenShift) & uint32(c.pixelFormat.GreenMax))
				color.B = uint16((rawPixel >> c.pixelFormat.BlueShift) & uint32(c.pixelFormat.BlueMax))
			} else {
				*color = c.colorMap[rawPixel]
			}
		}
	}

	return NewRawEncoding(colors), nil
}
