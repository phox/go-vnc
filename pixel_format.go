/*
pixel_format.go implements RFC 6143 §7.4 Pixel Format Data Structure.
See http://tools.ietf.org/html/rfc6143#section-7.4 for more info.
*/
package vnc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const pfSize = 16 // PixelFormat structure size.

// PixelFormat describes the way a pixel is formatted for a VNC connection.
type PixelFormat struct {
	BPP, Depth, BigEndian, TrueColor uint8
	RedMax, GreenMax, BlueMax        uint16 // 2^BPP-1
	RedShift, GreenShift, BlueShift  uint8
}

// NewPixelFormat returns a populated PixelFormat structure.
func NewPixelFormat() PixelFormat {
	return PixelFormat{16, 16, RFBTrue, RFBTrue, uint16(math.Exp2(16) - 1), uint16(math.Exp2(16) - 1), uint16(math.Exp2(16) - 1), 0, 0, 0}
}

// Bytes returns a slice of the contents of the PixelFormat structure. If there
// is an error creating the slice, an error will be returned.
func (pf PixelFormat) Bytes() ([]byte, error) {
	var buf bytes.Buffer

	// Validation checks.
	switch pf.BPP {
	case 8, 16, 32:
	default:
		return nil, NewVNCError(fmt.Sprintf("Invalid BPP value %v; must be 8, 16, or 32.", pf.BPP))
	}

	if pf.Depth < pf.BPP {
		return nil, NewVNCError(fmt.Sprintf("Invalid Depth value %v; must be >= BPP", pf.Depth))
	}
	switch pf.Depth {
	case 8, 16, 32:
	default:
		return nil, NewVNCError(fmt.Sprintf("Invalid Depth value %v; must be 8, 16, or 32.", pf.Depth))
	}

	// Create the slice of bytes
	if err := binary.Write(&buf, binary.BigEndian, pf); err != nil {
		return nil, err
	}

	pad := [3]byte{}
	if err := binary.Write(&buf, binary.BigEndian, pad); err != nil {
		return nil, err
	}

	// Padding values automatically set to 0 during slice conversion.
	return buf.Bytes()[:pfSize], nil
}

// Write populates the PixelFormat structure with data from the io.Reader. Any
// error encountered will be returned.
func (pf *PixelFormat) Write(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, pf); err != nil {
		return err
	}

	var pad [3]byte
	if err := binary.Read(r, binary.BigEndian, &pad); err != nil {
		return err
	}

	if pf.TrueColor != RFBFalse {
		pf.TrueColor = RFBTrue // Convert all non-zero values to our constant value.
	}

	return nil
}
