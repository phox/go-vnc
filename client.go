// Implementation of §7.5 Client-to-Server Messages.

package vnc

import (
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
)

const (
	setPixelFormatMsg = uint8(iota)
	_
	setEncodingsMsg
	framebufferUpdateRequestMsg
	keyEventMsg
	pointerEventMsg
	clientCutTextMsg
)

// SetPixelFormatMessage holds the wire format message.
type SetPixelFormatMessage struct {
	Msg uint8       // message-type
	_   [3]byte     // padding
	PF  PixelFormat // pixel-format
}

// SetPixelFormat sets the format in which pixel values should be sent
// in FramebufferUpdate messages from the server.
//
// See RFC 6143 Section 7.5.1
func (c *ClientConn) SetPixelFormat(pf PixelFormat) error {
	msg := SetPixelFormatMessage{
		Msg: setPixelFormatMsg,
		PF:  pf,
	}
	if err := c.send(msg); err != nil {
		return err
	}

	// Invalidate the color map.
	if pf.TrueColor == RFBFalse {
		c.colorMap = [256]Color{}
	}

	c.pixelFormat = pf
	return nil
}

// SetEncodingsMessage holds the wire format message, sans encoding-type field.
type SetEncodingsMessage struct {
	Msg     uint8   // message-type
	_       [1]byte // padding
	NumEncs uint16  // number-of-encodings
}

// SetEncodings sets the encoding types in which the pixel data can be sent
// from the server. After calling this method, the encs slice given should not
// be modified.
//
// See RFC 6143 Section 7.5.2
func (c *ClientConn) SetEncodings(encs Encodings) error {
	// Make sure RawEncoding is supported.
	haveRaw := false
	for _, v := range encs {
		if v.Type() == Raw {
			haveRaw = true
			break
		}
	}
	if !haveRaw {
		encs = append(encs, &RawEncoding{})
	}

	buf := NewBuffer(nil)

	// Prepare message.
	msg := SetEncodingsMessage{
		Msg:     setEncodingsMsg,
		NumEncs: uint16(len(encs)),
	}
	if err := buf.Write(msg); err != nil {
		return err
	}
	bytes, err := encs.Marshal()
	if err != nil {
		return err
	}
	buf.WriteBytes(bytes)

	// Send message.
	if err := c.send(buf.Bytes()); err != nil {
		return err
	}

	c.encodings = encs
	return nil
}

// FramebufferUpdateRequestMessage holds the wire format message.
type FramebufferUpdateRequestMessage struct {
	Msg           uint8  // message-type
	Inc           uint8  // incremental
	X, Y          uint16 // x-, y-position
	Width, Height uint16 // width, height
}

// Requests a framebuffer update from the server. There may be an indefinite
// time between the request and the actual framebuffer update being received.
//
// See RFC 6143 Section 7.5.3
func (c *ClientConn) FramebufferUpdateRequest(inc uint8, x, y, w, h uint16) error {
	msg := FramebufferUpdateRequestMessage{framebufferUpdateRequestMsg, inc, x, y, w, h}
	if err := c.send(&msg); err != nil {
		return err
	}
	return nil
}

// KeyEventMessage holds the wire format message.
type KeyEventMessage struct {
	Msg      uint8    // message-type
	DownFlag uint8    // down-flag
	_        [2]byte  // padding
	Key      keys.Key // key
}

const (
	PressKey   = true
	ReleaseKey = false
)

// KeyEvent indicates a key press or release and sends it to the server.
// The key is indicated using the X Window System "keysym" value. Use
// Google to find a reference of these values. To simulate a key press,
// you must send a key with both a down event, and a non-down event.
//
// See RFC 6143 Section 7.5.4.
func (c *ClientConn) KeyEvent(key keys.Key, down bool) error {
	if c.debug {
		log.Printf("KeyEvent(0x%04x, %v)", key, down)
	}

	downFlag := RFBFalse
	if down {
		downFlag = RFBTrue
	}
	msg := KeyEventMessage{keyEventMsg, downFlag, [2]byte{}, key}
	if err := c.send(msg); err != nil {
		return err
	}

	settleUI()
	return nil
}

// PointerEventMessage holds the wire format message.
type PointerEventMessage struct {
	Msg  uint8  // message-type
	Mask uint8  // button-mask
	X, Y uint16 // x-, y-position
}

// PointerEvent indicates that pointer movement or a pointer button
// press or release.
//
// The `button` is a bitwise mask of various Button values. When a button
// is set, it is pressed, when it is unset, it is released.
//
// See RFC 6143 Section 7.5.5
func (c *ClientConn) PointerEvent(button buttons.Button, x, y uint16) error {
	if c.debug {
		log.Printf("PointerEvent(%08b, %v, %v)", button, x, y)
	}

	msg := PointerEventMessage{pointerEventMsg, uint8(button), x, y}
	if err := c.send(msg); err != nil {
		return err
	}

	settleUI()
	return nil
}

// ClientCutTextMessage holds the wire format message, sans the text field.
type ClientCutTextMessage struct {
	Msg    uint8   // message-type
	_      [3]byte // padding
	Length uint32  // length
}

// ClientCutText tells the server that the client has new text in its cut buffer.
// The text string MUST only contain Latin-1 characters. This encoding
// is compatible with Go's native string format, but can only use up to
// unicode.MaxLatin1 values.
//
// See RFC 6143 Section 7.5.6
func (c *ClientConn) ClientCutText(text string) error {
	if c.debug {
		log.Printf("ClientCutText(%v)", text)
	}

	for _, char := range text {
		if char > unicode.MaxLatin1 {
			return NewVNCError(fmt.Sprintf("Character '%s' is not valid Latin-1", char))
		}
	}

	// Strip carriage-return (0x0d) chars.
	// From RFC: "Ends of lines are represented by the newline character (0x0a)
	// alone. No carriage-return (0x0d) is used."
	text = strings.Join(strings.Split(text, "\r"), "")

	msg := ClientCutTextMessage{
		Msg:    clientCutTextMsg,
		Length: uint32(len(text)),
	}
	if err := c.send(msg); err != nil {
		return err
	}
	if err := c.send([]byte(text)); err != nil {
		return err
	}

	settleUI()
	return nil
}
