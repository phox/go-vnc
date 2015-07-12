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

// RectangleMessage holds the wire format message.
type RectangleMessage struct {
	X, Y          uint16 // x-, y-position
	Width, Height uint16 // width, height
	Encoding      int32  // encoding-type
}

// Rectangle represents a rectangle of pixel data.
type Rectangle struct {
	X, Y          uint16
	Width, Height uint16
	Enc           Encoding
}

// FramebufferUpdate holds the wire format message.
type FramebufferUpdate struct {
	Msg     uint8       // message-type
	_       [1]byte     // padding
	NumRect uint16      // number-of-rectangles
	Rects   []Rectangle // rectangles
}

// FramebufferUpdateMessage holds the wire format message, sans rectangles.
type FramebufferUpdateMessage struct {
	Msg     uint8   // message-type
	_       [1]byte // padding
	NumRect uint16  // number-of-rectangles
}

func NewFramebufferUpdate(rects []Rectangle) *FramebufferUpdate {
	return &FramebufferUpdate{
		Msg:     FramebufferUpdateMsg,
		NumRect: uint16(len(rects)),
		Rects:   rects,
	}
}

func (m *FramebufferUpdate) Type() uint8 {
	return m.Msg
}

func (m *FramebufferUpdate) Read(c *ClientConn) (ServerMessage, error) {
	if c.debug {
		log.Print("FramebufferUpdate.Read()")
	}

	// Build the map of encodings supported
	encMap := make(map[int32]Encoding)
	for _, e := range c.Encodings() {
		encMap[e.Type()] = e
	}
	encMap[Raw] = NewRawEncoding([]Color{}) // Raw encoding support required.

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
		log.Printf("numRects:%v", numRects)
	}

	// Extract rectangles.
	rects := make([]Rectangle, numRects)
	for i := 0; i < int(numRects); i++ {
		var msg RectangleMessage
		if err := c.receive(&msg); err != nil {
			return nil, err
		}

		enc, ok := encMap[msg.Encoding]
		if !ok {
			return nil, fmt.Errorf("unsupported encoding type: %v; %v", msg.Encoding, msg)
		}

		rect := &rects[i]
		rect.X = msg.X
		rect.Y = msg.Y
		rect.Width = msg.Width
		rect.Height = msg.Height

		var err error
		rect.Enc, err = enc.Read(c, rect)
		if err != nil {
			return nil, fmt.Errorf("error reading encoding: %v", err)
		}
	}

	return NewFramebufferUpdate(rects), nil
}

// SetColorMapEntries is sent by the server to set values into
// the color map. This message will automatically update the color map
// for the associated connection, but contains the color change data
// if the consumer wants to read it.
//
// See RFC 6143 Section 7.6.2

// Color represents a single color in a color map.
type Color struct {
	R, G, B uint16
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
