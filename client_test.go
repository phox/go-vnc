package vnc

import (
	"encoding/binary"
	"testing"
)

func TestSetPixelFormat(t *testing.T) {
	t = nil
}

func TestSetEncodings(t *testing.T) {
	tests := []struct {
		encs     []Encoding
		encTypes []int32
	}{
		{[]Encoding{&RawEncoding{}}, []int32{0}},
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
			t.Errorf("FramebufferUpdateRequest() X incorrect; got = %v, want = %v", req.Width, tt.w)
		}
		if req.Height != tt.h {
			t.Errorf("FramebufferUpdateRequest() X incorrect; got = %v, want = %v", req.Height, tt.h)
		}
	}
}

func TestKeyEvent(t *testing.T) {
	t = nil
}

func TestPointerEvent(t *testing.T) {
	t = nil
}

func TestClientCutText(t *testing.T) {
	t = nil
}
