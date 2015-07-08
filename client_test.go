package vnc

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func TestSetPixelFormat(t *testing.T) {
	t = nil
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
		err := conn.SetEncodings(tt.encs)
		if err != nil {
			t.Fatalf("SendEncodingsMessage() unexpected error %v", err)
		}

		// Read the request.
		req := SetEncodingsMessage{}
		if err := binary.Read(conn.c, binary.BigEndian, &req.Msg); err != nil {
			t.Fatal(err)
		}
		for i := range req.Pad {
			if err := binary.Read(conn.c, binary.BigEndian, &req.Pad[i]); err != nil {
				t.Fatal(err)
			}
		}
		if err := binary.Read(conn.c, binary.BigEndian, &req.NumEncs); err != nil {
			t.Fatal(err)
		}
		var encs []int32 // Can't use the request struct.
		for i := 0; i < len(tt.encs); i++ {
			var enc int32
			if err := binary.Read(conn.c, binary.BigEndian, &enc); err != nil {
				t.Fatal(err)
			}
			encs = append(encs, enc)
		}

		// Validate the request.
		if req.Msg != setEncodingsMsg {
			t.Errorf("SetEncodings() incorrect message-type; got = %v, want = %v", req.Msg, setEncodingsMsg)
		}
		if req.NumEncs != uint16(len(tt.encs)) {
			t.Errorf("SetEncodings() incorrect number of encodings; got = %v, want = %v", req.NumEncs, len(tt.encs))
		}
		for i := 0; i < len(tt.encs); i++ {
			if encs[i] != tt.encs[i].Type() {
				t.Errorf("SetEncodings() incorrect encoding [%v}]; got = %v, want = %v", i, req.Encs[i], tt.encs[i].Type())
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
			t.Fatalf("FramebufferUpdateRequest() unexpected error %v", err)
		}

		// Validate the request.
		req := FramebufferUpdateRequestMessage{}
		if err := binary.Read(conn.c, binary.BigEndian, &req); err != nil {
			t.Fatal(err)
		}
		if req.Msg != framebufferUpdateRequestMsg {
			t.Errorf("FramebufferUpdateRequest() incorrect message-type; got = %v, want = %v", req.Msg, framebufferUpdateRequestMsg)
		}
		if req.Inc != tt.inc {
			t.Errorf("FramebufferUpdateRequest() incremental incorrect; got = %v, want = %v", req.Inc, tt.inc)
		}
		if req.X != tt.x {
			t.Errorf("FramebufferUpdateRequest() X incorrect; got = %v, want = %v", req.X, tt.x)
		}
		if req.Y != tt.y {
			t.Errorf("FramebufferUpdateRequest() Y incorrect; got = %v, want = %v", req.Y, tt.y)
		}
		if req.Width != tt.w {
			t.Errorf("FramebufferUpdateRequest() Width incorrect; got = %v, want = %v", req.Width, tt.w)
		}
		if req.Height != tt.h {
			t.Errorf("FramebufferUpdateRequest() Height incorrect; got = %v, want = %v", req.Height, tt.h)
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

func TestKeyEvent(t *testing.T) {}

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

	// Give the mouse some time to "settle" after moving.
	time.Sleep(10 * time.Millisecond)

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

	for _, tt := range tests {
		mockConn.Reset()

		// Send request.
		err := conn.PointerEvent(tt.mask, tt.x, tt.y)
		if err != nil {
			t.Fatalf("PointerEvent() unexpected error %v", err)
		}

		// Validate the request.
		req := PointerEventMessage{}
		if err := binary.Read(conn.c, binary.BigEndian, &req); err != nil {
			t.Fatal(err)
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

func TestClientCutText(t *testing.T) {}
