/*
client.go implements RFC 6143 §7.5 Client-to-Server Messages.
See http://tools.ietf.org/html/rfc6143#section-7.5 for more info.
*/
package vnc

import (
	"fmt"
	"strings"
	"unicode"
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
func (c *ClientConn) SetEncodings(e []Encoding) error {
	// Prepare message.
	msg := SetEncodingsMessage{
		Msg:     setEncodingsMsg,
		NumEncs: uint16(len(e)),
	}
	var encs []int32 // encoding-type

	for _, v := range e {
		encs = append(encs, int32(v.Type()))
	}

	// Send message.
	if err := c.send(msg); err != nil {
		return err
	}
	if err := c.sendN(encs); err != nil {
		return err
	}

	c.encodings = e
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
	Msg      uint8   // message-type
	DownFlag uint8   // down-flag
	_        [2]byte // padding
	Key      uint32  // key
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
func (c *ClientConn) KeyEvent(keysym uint32, down bool) error {
	var downFlag uint8 = RFBFalse
	if down {
		downFlag = RFBTrue
	}

	msg := KeyEventMessage{keyEventMsg, downFlag, [2]byte{}, keysym}
	if err := c.send(msg); err != nil {
		return err
	}

	settleUI()
	return nil
}

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

// PointerEventMessage holds the wire format message.
type PointerEventMessage struct {
	Msg  uint8  // message-type
	Mask uint8  // button-mask
	X, Y uint16 // x-, y-position
}

// PointerEvent indicates that pointer movement or a pointer button
// press or release.
//
// The mask is a bitwise mask of various ButtonMask values. When a button
// is set, it is pressed, when it is unset, it is released.
//
// See RFC 6143 Section 7.5.5
func (c *ClientConn) PointerEvent(mask ButtonMask, x, y uint16) error {
	msg := PointerEventMessage{pointerEventMsg, uint8(mask), x, y}
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
	if err := c.sendN([]byte(text)); err != nil {
		return err
	}

	settleUI()
	return nil
}

//
// Constants to use for KeyEvents PointerEvents.
//

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
	KeyHome = iota + 0xff50
	KeyLeft
	KeyUp
	KeyRight
	KeyDown
	KeyPageUp
	KeyPageDown
	KeyEnd
	KeyInsert = 0xff63
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
