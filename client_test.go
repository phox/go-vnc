package vnc

import (
	"encoding/binary"
	"testing"
)

func TestSetPixelFormat(t *testing.T) {
	t = nil
}

func TestSetEncodings(t *testing.T) {
	t = nil
}

func TestFramebufferUpdateRequest(t *testing.T) {
	tests := []struct {
		inc        uint8
		x, y, w, h uint16
	}{
		{RFBFalse, 10, 20, 30, 40},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		err := conn.FramebufferUpdateRequest(tt.inc, tt.x, tt.y, tt.w, tt.h)
		if err != nil {
			t.Fatalf("FramebufferUpdateRequest() unexpected error %v", err)
		}

		// Validate client request.
		req := FramebufferUpdateRequestType{}
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
