package vnc

import (
	"fmt"
	"math"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/kward/go-vnc/go/operators"
	"golang.org/x/net/context"
)

func TestSetPixelFormat(t *testing.T) {
	max := uint16(math.Exp2(16))
	pf := PixelFormat{
		BPP:       16,
		Depth:     16,
		BigEndian: RFBTrue,
		TrueColor: RFBTrue,
		RedMax:    max,
		GreenMax:  max,
		BlueMax:   max,
	}

	tests := []struct {
		pf  PixelFormat
		msg SetPixelFormatMessage
	}{
		{PixelFormat{},
			SetPixelFormatMessage{
				Msg: setPixelFormatMsg,
			}},
		{NewPixelFormat(),
			SetPixelFormatMessage{
				Msg: setPixelFormatMsg,
				PF:  pf,
			}},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		if err := conn.SetPixelFormat(tt.pf); err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		// Read the request.
		req := SetPixelFormatMessage{}
		if err := conn.receive(&req); err != nil {
			t.Error(err)
			continue
		}
		if got, want := req.Msg, setPixelFormatMsg; got != want {
			t.Errorf("incorrect message-type; got = %v, want = %v", got, want)
		}
		if got, want := req.PF.BPP, tt.msg.PF.BPP; got != want {
			t.Errorf("incorrect pixel-format bits-per-pixel; got = %v, want = %v", got, want)
		}
		if got, want := req.PF.BlueShift, tt.msg.PF.BlueShift; got != want {
			t.Errorf("incorrect pixel-format blue-shift; got = %v, want = %v", got, want)
		}
	}
}

func TestSetEncodings(t *testing.T) {
	tests := []struct {
		encs     []Encoding
		encTypes []int32
	}{
		{[]Encoding{NewRawEncoding([]Color{})}, []int32{0}},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		if err := conn.SetEncodings(tt.encs); err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		// Read the request.
		req := SetEncodingsMessage{}
		if err := conn.receive(&req.Msg); err != nil {
			t.Error(err)
			continue
		}
		for i := range req.Pad {
			if err := conn.receive(&req.Pad[i]); err != nil {
				t.Error(err)
				continue
			}
		}
		if err := conn.receive(&req.NumEncs); err != nil {
			t.Error(err)
			continue
		}
		var encs []int32 // Can't use the request struct.
		for i := 0; i < len(tt.encs); i++ {
			var enc int32
			if err := conn.receive(&enc); err != nil {
				t.Error(err)
				continue
			}
			encs = append(encs, enc)
		}

		// Validate the request.
		if req.Msg != setEncodingsMsg {
			t.Errorf("incorrect message-type; got = %v, want = %v", req.Msg, setEncodingsMsg)
		}
		if req.NumEncs != uint16(len(tt.encs)) {
			t.Errorf("incorrect number-of-encodings; got = %v, want = %v", req.NumEncs, len(tt.encs))
		}
		for i := 0; i < len(tt.encs); i++ {
			if encs[i] != tt.encs[i].Type() {
				t.Errorf("incorrect encoding-type [%v]; got = %v, want = %v", i, req.Encs[i], tt.encs[i].Type())
			}
		}
	}
}

func TestFramebufferUpdateRequest(t *testing.T) {
	tests := []struct {
		inc        uint8
		x, y, w, h uint16
	}{
		{RFBFalse, 10, 20, 30, 40},
		{RFBTrue, 11, 21, 31, 41},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		err := conn.FramebufferUpdateRequest(tt.inc, tt.x, tt.y, tt.w, tt.h)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		// Validate the request.
		req := FramebufferUpdateRequestMessage{}
		if err := conn.receive(&req); err != nil {
			t.Error(err)
			continue
		}
		if req.Msg != framebufferUpdateRequestMsg {
			t.Errorf("incorrect message-type; got = %v, want = %v", req.Msg, framebufferUpdateRequestMsg)
		}
		if req.Inc != tt.inc {
			t.Errorf("incorrect incremental; got = %v, want = %v", req.Inc, tt.inc)
		}
		if req.X != tt.x {
			t.Errorf("incorrect x-position; got = %v, want = %v", req.X, tt.x)
		}
		if req.Y != tt.y {
			t.Errorf("incorrect y-position; got = %v, want = %v", req.Y, tt.y)
		}
		if req.Width != tt.w {
			t.Errorf("incorrect width; got = %v, want = %v", req.Width, tt.w)
		}
		if req.Height != tt.h {
			t.Errorf("incorrect height; got = %v, want = %v", req.Height, tt.h)
		}
	}
}

func ExampleClientConn_KeyEvent() {
	// Establish TCP connection.
	nc, err := net.DialTimeout("tcp", "127.0.0.1:5900", 10*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to host: %v\n", err))
	}

	// Negotiate VNC connection.
	vc, err := Connect(context.Background(), nc, NewClientConfig("somepass"))
	if err != nil {
		panic(fmt.Sprintf("Could not negotiate a VNC connection: %v\n", err))
	}

	// Press and release the return key.
	vc.KeyEvent(KeyReturn, true)
	vc.KeyEvent(KeyReturn, false)

	// Close VNC connection.
	vc.Close()
}

func TestKeyEvent(t *testing.T) {
	tests := []struct {
		key  uint32
		down bool
	}{
		{Key0, PressKey},
		{Key1, ReleaseKey},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	SetSettle(0) // Disable UI settling for tests.
	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		err := conn.KeyEvent(tt.key, tt.down)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		// Validate the request.
		req := KeyEventMessage{}
		if err := conn.receive(&req); err != nil {
			t.Error(err)
			continue
		}
		var down bool
		switch req.DownFlag {
		case RFBTrue:
			down = PressKey
		case RFBFalse:
			down = ReleaseKey
		}

		if got, want := req.Msg, keyEventMsg; got != want {
			t.Errorf("incorrect message-type; got = %v, want = %v", got, want)
		}
		if got, want := down, tt.down; got != want {
			t.Errorf("incorrect down-flag; got = %v, want = %v", got, want)
		}
		if got, want := req.Key, tt.key; got != want {
			t.Errorf("incorrect key; got = %v, want = %v", got, want)
		}
	}
}

func ExampleClientConn_PointerEvent() {
	// Establish TCP connection.
	nc, err := net.DialTimeout("tcp", "127.0.0.1:5900", 10*time.Second)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to host: %v\n", err))
	}

	// Negotiate VNC connection.
	vc, err := Connect(context.Background(), nc, NewClientConfig("somepass"))
	if err != nil {
		panic(fmt.Sprintf("Could not negotiate a VNC connection: %v\n", err))
	}

	// Move mouse to x=100, y=200.
	x, y := uint16(100), uint16(200)
	vc.PointerEvent(ButtonNone, x, y)

	// Click and release the left mouse button.
	vc.PointerEvent(ButtonLeft, x, y)
	vc.PointerEvent(ButtonNone, x, y)

	// Close connection.
	vc.Close()
}

func TestPointerEvent(t *testing.T) {
	tests := []struct {
		mask ButtonMask
		x, y uint16
	}{
		{ButtonNone, 0, 0},
		{ButtonLeft | ButtonRight, 123, 456},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	SetSettle(0) // Disable UI settling for tests.
	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		err := conn.PointerEvent(tt.mask, tt.x, tt.y)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}

		// Validate the request.
		req := PointerEventMessage{}
		if err := conn.receive(&req); err != nil {
			t.Error(err)
			continue
		}
		if got, want := req.Msg, pointerEventMsg; got != want {
			t.Errorf("incorrect message-type; got = %v, want = %v", got, want)
		}
		if got, want := req.Mask, uint8(tt.mask); got != want {
			t.Errorf("incorrect button-mask; got = %v, want = %v", got, want)
		}
		if got, want := req.X, tt.x; got != want {
			t.Errorf("incorrect x-position; got = %v, want = %v", got, want)
		}
		if got, want := req.Y, tt.y; got != want {
			t.Errorf("incorrect y-position; got = %v, want = %v", got, want)
		}
	}
}

func TestClientCutText(t *testing.T) {
	tests := []struct {
		text string
		sent []byte
		ok   bool
	}{
		{"abc123", []byte("abc123"), true},
		{"foo\r\nbar", []byte("foo\nbar"), true},
		{"", []byte{}, true},
		{"ɹɐqooɟ", []byte{}, false},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	SetSettle(0) // Disable UI settling for tests.
	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		err := conn.ClientCutText(tt.text)
		if err == nil && !tt.ok {
			t.Errorf("expected error")
		}
		if err != nil {
			if verr, ok := err.(*VNCError); !ok {
				t.Errorf("unexpected %v error: %v", reflect.TypeOf(err), verr)
			}
		}
		if !tt.ok {
			continue
		}

		// Validate the request.
		req := ClientCutTextMessage{}
		if err := conn.receive(&req); err != nil {
			t.Error(err)
			continue
		}
		if got, want := req.Msg, clientCutTextMsg; got != want {
			t.Errorf("incorrect message-type; got = %v, want = %v", got, want)
		}
		if got, want := req.Length, uint32(len(tt.sent)); got != want {
			t.Errorf("incorrect length; got = %v, want = %v", got, want)
		}
		var text []byte
		if err := conn.receiveN(&text, int(req.Length)); err != nil {
			t.Error(err)
			continue
		}
		if got, want := text, tt.sent; !operators.EqualSlicesOfByte(got, want) {
			t.Errorf("incorrect text; got = %v, want = %v", got, want)
		}
	}
}
