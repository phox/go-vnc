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
	BPP        uint8
	Depth      uint8
	BigEndian  uint8
	TrueColor  uint8
	RedMax     uint16 // 2^BPP-1
	GreenMax   uint16 // 2^BPP-1
	BlueMax    uint16 // 2^BPP-1
	RedShift   uint8
	GreenShift uint8
	BlueShift  uint8
	padding    [3]byte
}

// NewPixelFormat returns a populated PixelFormat structure.
func NewPixelFormat() PixelFormat {
	return PixelFormat{16, 16, RFBTrue, RFBTrue, uint16(math.Exp2(16) - 1), uint16(math.Exp2(16) - 1), uint16(math.Exp2(16) - 1), 0, 0, 0, [3]byte{}}
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
	if err := binary.Write(&buf, binary.BigEndian, pf.BPP); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, pf.Depth); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, pf.BigEndian); err != nil {
		return nil, err
	}
	if err := binary.Write(&buf, binary.BigEndian, pf.TrueColor); err != nil {
		return nil, err
	}

	// If TrueColor is true, then populate the structure with the color values.
	if pf.TrueColor == RFBTrue {
		if err := binary.Write(&buf, binary.BigEndian, pf.RedMax); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, pf.GreenMax); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, pf.BlueMax); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, pf.RedShift); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, pf.GreenShift); err != nil {
			return nil, err
		}
		if err := binary.Write(&buf, binary.BigEndian, pf.BlueShift); err != nil {
			return nil, err
		}
	}

	// Padding values automatically set to 0 during slice conversion.
	return buf.Bytes()[:pfSize], nil
}

// Write populates the PixelFormat structure with data from the io.Reader. Any
// error encountered will be returned.
func (pf *PixelFormat) Write(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, &pf.BPP); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.Depth); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.BigEndian); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.TrueColor); err != nil {
		return err
	}
	if pf.TrueColor != RFBFalse {
		pf.TrueColor = RFBTrue // Convert all non-zero values to our constant value.
	}

	if err := binary.Read(r, binary.BigEndian, &pf.RedMax); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.GreenMax); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.BlueMax); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.RedShift); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.GreenShift); err != nil {
		return err
	}
	if err := binary.Read(r, binary.BigEndian, &pf.BlueShift); err != nil {
		return err
	}

	var padding [3]uint8
	if err := binary.Read(r, binary.BigEndian, &padding); err != nil {
		return err
	}

	return nil
}
