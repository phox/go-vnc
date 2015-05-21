/*
Initialization implements RFC 6143 §7.5 Client-to-Server Messages.

See http://tools.ietf.org/html/rfc6143#section-7.5 for more info.

Sample usage:

// Move mouse to x=100, y=200.
x, y := 100, 200
conn.PointerEvent(vnc.ButtonNone, x, y)
// Give mouse some time to "settle."
time.Sleep(10*time.Millisecond)
// Left click.
conn.PointerEvent(vnc.ButtonLeft, x, y)
conn.PointerEvent(vnc.ButtonNone, x, y)

// Press return key
conn.KeyEvent(vnc.KeyReturn, true)
// Release the key.
conn.KeyEvent(vnc.KeyReturn, false)
*/
package vnc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode"
)

const (
	setPixelFormatMsg = iota
	_
	setEncodingsMsg
	framebufferUpdateRequestMsg
	keyEventMsg
	pointerEventMsg
	clientCutTextMsg
)

// SetPixelFormat sets the format in which pixel values should be sent
// in FramebufferUpdate messages from the server.
//
// See RFC 6143 Section 7.5.1
func (c *ClientConn) SetPixelFormat(format *PixelFormat) error {
	var keyEvent [20]byte
	keyEvent[0] = 0

	pfBytes, err := writePixelFormat(format)
	if err != nil {
		return err
	}

	// Copy the pixel format bytes into the proper slice location
	copy(keyEvent[4:], pfBytes)

	// Send the data down the connection
	if _, err := c.c.Write(keyEvent[:]); err != nil {
		return err
	}

	// Reset the color map as according to RFC.
	var newColorMap [256]Color
	c.ColorMap = newColorMap

	return nil
}

// SetEncodings sets the encoding types in which the pixel data can
// be sent from the server. After calling this method, the encs slice
// given should not be modified.
//
// See RFC 6143 Section 7.5.2
func (c *ClientConn) SetEncodings(encs []Encoding) error {
	data := make([]interface{}, 3+len(encs))
	data[0] = uint8(setEncodingsMsg)
	data[1] = uint8(0)
	data[2] = uint16(len(encs))

	for i, enc := range encs {
		data[3+i] = int32(enc.Type())
	}

	var buf bytes.Buffer
	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	dataLength := 4 + (4 * len(encs))
	if _, err := c.c.Write(buf.Bytes()[0:dataLength]); err != nil {
		return err
	}

	c.Encs = encs

	return nil
}

// Requests a framebuffer update from the server. There may be an indefinite
// time between the request and the actual framebuffer update being
// received.
//
// See RFC 6143 Section 7.5.3
func (c *ClientConn) FramebufferUpdateRequest(incremental bool, x, y, width, height uint16) error {
	var buf bytes.Buffer
	var incrementalByte uint8

	if incremental {
		incrementalByte = 1
	}

	data := []interface{}{
		uint8(framebufferUpdateMsg),
		incrementalByte,
		x, y, width, height,
	}

	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}
	if _, err := c.c.Write(buf.Bytes()[0:10]); err != nil {
		return err
	}

	return nil
}

// KeyEvent indicates a key press or release and sends it to the server.
// The key is indicated using the X Window System "keysym" value. Use
// Google to find a reference of these values. To simulate a key press,
// you must send a key with both a down event, and a non-down event.
//
// See 7.5.4.
func (c *ClientConn) KeyEvent(keysym uint32, down bool) error {
	var downFlag uint8 = 0
	if down {
		downFlag = 1
	}

	data := []interface{}{
		uint8(keyEventMsg),
		downFlag,
		uint8(0),
		uint8(0),
		keysym,
	}

	for _, val := range data {
		if err := binary.Write(c.c, binary.BigEndian, val); err != nil {
			return err
		}
	}

	return nil
}

// PointerEvent indicates that pointer movement or a pointer button
// press or release.
//
// The mask is a bitwise mask of various ButtonMask values. When a button
// is set, it is pressed, when it is unset, it is released.
//
// See RFC 6143 Section 7.5.5
func (c *ClientConn) PointerEvent(mask ButtonMask, x, y uint16) error {
	var buf bytes.Buffer

	data := []interface{}{
		uint8(pointerEventMsg),
		uint8(mask),
		x,
		y,
	}

	for _, val := range data {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	if _, err := c.c.Write(buf.Bytes()[0:6]); err != nil {
		return err
	}

	return nil
}

// ClientCutText tells the server that the client has new text in its cut buffer.
// The text string MUST only contain Latin-1 characters. This encoding
// is compatible with Go's native string format, but can only use up to
// unicode.MaxLatin values.
//
// See RFC 6143 Section 7.5.6
func (c *ClientConn) ClientCutText(text string) error {
	var buf bytes.Buffer

	// This is the fixed size data we'll send
	fixedData := []interface{}{
		uint8(clientCutTextMsg),
		uint8(0),
		uint8(0),
		uint8(0),
		uint32(len(text)),
	}

	for _, val := range fixedData {
		if err := binary.Write(&buf, binary.BigEndian, val); err != nil {
			return err
		}
	}

	for _, char := range text {
		if char > unicode.MaxLatin1 {
			return fmt.Errorf("Character '%s' is not valid Latin-1", char)
		}

		if err := binary.Write(&buf, binary.BigEndian, uint8(char)); err != nil {
			return err
		}
	}

	dataLength := 8 + len(text)
	if _, err := c.c.Write(buf.Bytes()[0:dataLength]); err != nil {
		return err
	}

	return nil
}

//
// Constants to use for KeyEvents PointerEvents.
//

// ButtonMask represents a mask of pointer presses/releases.
type ButtonMask uint8

// All available button mask components.
const (
	ButtonLeft ButtonMask = 1 << iota
	ButtonMiddle
	ButtonRight
	Button4
	Button5
	Button6
	Button7
	Button8
	ButtonNone = ButtonMask(0)
)

// Latin 1 (byte 3 = 0)
// ISO/IEC 8859-1 = Unicode U+0020..U+00FF
const (
	KeySpace = iota + 0x0020
	KeyExclam
	KeyQuoteDbl
	KeyNumberSign
	KeyDollar
	KeyPercent
	KeyAmpersand
	KeyApostrophe
	KeyParenLeft
	KeyParenRight
	KeyAsterisk
	KeyPlus
	KeyComma
	KeyMinus
	KeyPeriod
	KeySlash
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeyColon
	KeySemicolon
	KeyLess
	KeyEqual
	KeyGreater
	KeyQuestion
	KeyAt
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyBracketLeft
	KeyBackslash
	KeyBracketRight
	KeyAsciiCircum
	KeyUnderscore
	KeyGrave
	Keya
	Keyb
	Keyc
	Keyd
	Keye
	Keyf
	Keyg
	Keyh
	Keyi
	Keyj
	Keyk
	Keyl
	Keym
	Keyn
	Keyo
	Keyp
	Keyq
	Keyr
	Keys
	Keyt
	Keyu
	Keyv
	Keyw
	Keyx
	Keyy
	Keyz
	KeyBraceLeft
	KeyBar
	KeyBraceRight
	KeyAsciiTilde
)
const (
	KeyBackspace = iota + 0xff08
	KeyTab
	KeyLinefeed
	KeyClear
	_
	KeyReturn
)
const (
	KeyPause      = 0xff13
	KeyScrollLock = 0xff14
	KeySysReq     = 0xff15
	KeyEscape     = 0xff1b
)
const (
	KeyF1 = iota + 0xffbe
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
)
const (
	KeyShiftLeft = iota + 0xffe1
	KeyShiftRight
	KeyControlLeft
	KeyControlRight
	KeyCapsLock
	_
	_
	_
	KeyAltLeft
	KeyAltRight

	KeyDelete = 0xffff
)
