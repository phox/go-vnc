// Implementation of RFC 6143 §7.6 Server-to-Client Messages.

package vnc

import (
	"fmt"
	"log"
)

const (
	FramebufferUpdateMsg = uint8(iota)
	SetColorMapEntriesMsg
	BellMsg
	ServerCutTextMsg
)

// A ServerMessage implements a message sent from the server to the client.
type ServerMessage interface {
	// The type of the message that is sent down on the wire.
	Type() uint8

	// Read reads the contents of the message from the reader. At the point
	// this is called, the message type has already been read from the reader.
	// This should return a new ServerMessage that is the appropriate type.
	Read(*ClientConn) (ServerMessage, error)
}

// Rectangle wire format message.
type rectangleMessage struct {
	X, Y uint16 // x-, y-position
	W, H uint16 // width, height
	E    int32  // encoding-type
}

// Rectangle represents a rectangle of pixel data.
type Rectangle struct {
	X, Y          uint16
	Width, Height uint16
	Enc           Encoding
	encodable     EncodableFunc
}

func NewRectangle(c *ClientConn) *Rectangle {
	return &Rectangle{encodable: c.Encodable}
}

// Read a rectangle from a connection.
func (r *Rectangle) Read(c *ClientConn) error {
	var msg rectangleMessage
	if err := c.receive(&msg); err != nil {
		return err
	}
	r.X, r.Y, r.Width, r.Height = msg.X, msg.Y, msg.W, msg.H

	enc, ok := r.encodable(msg.E)
	if !ok {
		return fmt.Errorf("unsupported encoding type: %v", msg.E)
	}

	var err error
	r.Enc, err = enc.Read(c, r)
	if err != nil {
		return fmt.Errorf("error reading rectangle encoding: %v", err)
	}

	return nil
}

func (r *Rectangle) Marshal() ([]byte, error) {
	buf := NewBuffer(nil)

	var msg rectangleMessage
	msg.X, msg.Y, msg.W, msg.H = r.X, r.Y, r.Width, r.Height
	msg.E = r.Enc.Type()
	if err := buf.Write(msg); err != nil {
		return nil, err
	}

	bytes, err := r.Enc.Marshal()
	if err != nil {
		return nil, err
	}
	buf.WriteBytes(bytes)

	return buf.Bytes(), nil
}

func (r *Rectangle) Unmarshal(data []byte) error {
	buf := NewBuffer(data)

	var msg rectangleMessage
	if err := buf.Read(&msg); err != nil {
		return err
	}
	r.X, r.Y, r.Width, r.Height = msg.X, msg.Y, msg.W, msg.H

	switch msg.E {
	case Raw:
		r.Enc = &RawEncoding{}
	default:
		return fmt.Errorf("unable to unmarshal encoding %v", msg.E)
	}
	return nil
}

func (r *Rectangle) Area() int { return int(r.Width) * int(r.Height) }

func (r *Rectangle) DebugPrint() {
	log.Printf("Rectangle: x: %v y: %v, w: %v, h: %v, enc: %v", r.X, r.Y, r.Width, r.Height, r.Enc)
}

type EncodableFunc func(enc int32) (Encoding, bool)

func (c *ClientConn) Encodable(enc int32) (Encoding, bool) {
	for _, e := range c.encodings {
		if e.Type() == enc {
			return e, true
		}
	}
	return nil, false
}

// FramebufferUpdate holds the wire format message.
type FramebufferUpdate struct {
	NumRect uint16      // number-of-rectangles
	Rects   []Rectangle // rectangles
}

func NewFramebufferUpdate(rects []Rectangle) *FramebufferUpdate {
	return &FramebufferUpdate{
		NumRect: uint16(len(rects)),
		Rects:   rects,
	}
}

func (m *FramebufferUpdate) Type() uint8 {
	return FramebufferUpdateMsg
}

func (m *FramebufferUpdate) Read(c *ClientConn) (ServerMessage, error) {
	if c.debug {
		log.Print("FramebufferUpdate.Read()")
	}

	// Build the map of supported encodings.
	// encs := make(map[int32]Encoding)
	// for _, e := range c.Encodings() {
	// 	encs[e.Type()] = e
	// }
	// encs[Raw] = &RawEncoding{} // Raw encoding support required.

	// Read packet.
	var pad [1]byte
	if err := c.receive(&pad); err != nil {
		return nil, err
	}
	var numRects uint16
	if err := c.receive(&numRects); err != nil {
		return nil, err
	}
	if c.debug {
		log.Printf("FramebufferUpdate.Read() numRects: %v", numRects)
	}

	// Extract rectangles.
	rects := make([]Rectangle, numRects)
	for i := 0; i < int(numRects); i++ {
		rect := NewRectangle(c)
		if err := rect.Read(c); err != nil {
			return nil, err
		}
		rects[i] = *rect
	}

	return NewFramebufferUpdate(rects), nil
}

func (m *FramebufferUpdate) Marshal() ([]byte, error) {
	buf := NewBuffer(nil)

	msg := struct {
		msg      uint8   // message-type
		_        [1]byte // padding
		numRects uint16  // number-of-rectangles
	}{
		msg:      FramebufferUpdateMsg,
		numRects: m.NumRect,
	}
	if err := buf.Write(msg); err != nil {
		return nil, err
	}
	for _, rect := range m.Rects {
		bytes, err := rect.Marshal()
		if err != nil {
			return nil, err
		}
		buf.WriteBytes(bytes)
	}

	return buf.Bytes(), nil
}

// SetColorMapEntries is sent by the server to set values into
// the color map. This message will automatically update the color map
// for the associated connection, but contains the color change data
// if the consumer wants to read it.
//
// See RFC 6143 Section 7.6.2

// Color represents a single color in a color map.
type Color struct {
	pf      *PixelFormat
	cm      *ColorMap
	cmIndex uint32 // Only valid if pf.TrueColor is false.
	R, G, B uint16
}

type ColorMap [256]Color

func NewColor(pf *PixelFormat, cm *ColorMap) *Color {
	return &Color{pf: pf, cm: cm}
}

func (c *Color) Marshal() ([]byte, error) {
	order := c.pf.order()

	var pixel uint32
	if c.pf.TrueColor == RFBTrue {
		pixel = uint32(c.R) << c.pf.RedShift
		pixel |= uint32(c.G) << c.pf.GreenShift
		pixel |= uint32(c.B) << c.pf.BlueShift
	} else {
		pixel = c.cmIndex
	}

	var bytes []byte
	switch c.pf.BPP {
	case 8:
		bytes = make([]byte, 1)
		bytes[0] = byte(pixel)
	case 16:
		bytes = make([]byte, 2)
		order.PutUint16(bytes, uint16(pixel))
	case 32:
		bytes = make([]byte, 4)
		order.PutUint32(bytes, pixel)
	}

	return bytes, nil
}

func (c *Color) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return NewVNCError(fmt.Sprint("Could not unmarshal empty data slice"))
	}
	order := c.pf.order()

	var pixel uint32
	switch c.pf.BPP {
	case 8:
		pixel = uint32(data[0])
	case 16:
		pixel = uint32(order.Uint16(data))
	case 32:
		pixel = order.Uint32(data)
	}

	if c.pf.TrueColor == RFBTrue {
		c.R = uint16((pixel >> c.pf.RedShift) & uint32(c.pf.RedMax))
		c.G = uint16((pixel >> c.pf.GreenShift) & uint32(c.pf.GreenMax))
		c.B = uint16((pixel >> c.pf.BlueShift) & uint32(c.pf.BlueMax))
	} else {
		*c = c.cm[pixel]
		c.cmIndex = pixel
	}

	return nil
}

type SetColorMapEntries struct {
	FirstColor uint16
	Colors     []Color
}

func (*SetColorMapEntries) Type() uint8 {
	return SetColorMapEntriesMsg
}

func (*SetColorMapEntries) Read(c *ClientConn) (ServerMessage, error) {
	if c.debug {
		log.Print("SetColorMapEntries.Read()")
	}

	// Read off the padding
	var padding [1]byte
	if err := c.receive(&padding); err != nil {
		return nil, err
	}

	var result SetColorMapEntries
	if err := c.receive(&result.FirstColor); err != nil {
		return nil, err
	}

	var numColors uint16
	if err := c.receive(&numColors); err != nil {
		return nil, err
	}

	result.Colors = make([]Color, numColors)
	for i := uint16(0); i < numColors; i++ {

		color := &result.Colors[i]
		if err := c.receive(&color); err != nil {
			return nil, err
		}

		// Update the connection's color map
		c.colorMap[result.FirstColor+i] = *color
	}

	return &result, nil
}

// Bell signals that an audible bell should be made on the client.
//
// See RFC 6143 Section 7.6.3
type Bell struct{}

func (*Bell) Type() uint8 {
	return BellMsg
}

func (*Bell) Read(c *ClientConn) (ServerMessage, error) {
	if c.debug {
		log.Print("Bell.Read()")
	}

	return new(Bell), nil
}

// ServerCutText indicates the server has new text in the cut buffer.
//
// See RFC 6143 Section 7.6.4
type ServerCutText struct {
	Text string
}

func (*ServerCutText) Type() uint8 {
	return ServerCutTextMsg
}

func (*ServerCutText) Read(c *ClientConn) (ServerMessage, error) {
	if c.debug {
		log.Print("ServerCutText.Read()")
	}

	// Read off the padding
	var padding [1]byte
	if err := c.receive(&padding); err != nil {
		return nil, err
	}

	var textLength uint32
	if err := c.receive(&textLength); err != nil {
		return nil, err
	}

	textBytes := make([]uint8, textLength)
	if err := c.receive(&textBytes); err != nil {
		return nil, err
	}

	return &ServerCutText{string(textBytes)}, nil
}
