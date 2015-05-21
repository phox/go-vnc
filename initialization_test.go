package vnc

import (
	"encoding/binary"
	"testing"
)

func TestClientInit(t *testing.T) {
	tests := []struct {
		exclusive bool
		shared    uint8
	}{
		{true, 0},
		{false, 1},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()
		conn.config.Exclusive = tt.exclusive

		// Validate client response.
		err := conn.clientInit()
		if err != nil {
			t.Fatalf("clientInit() unexpected error %v", err)
		}
		var shared uint8
		err = binary.Read(conn.c, binary.BigEndian, &shared)
		if shared != tt.shared {
			t.Errorf("clientInit() shared: got = %v, want = %v", shared, tt.shared)
		}
	}
}

func TestServerInit(t *testing.T) {
	const (
		none = iota
		fbw
		fbh
		pf
		dn
	)
	tests := []struct {
		eof               int
		fbWidth, fbHeight uint16
		pixelFormat       PixelFormat
		desktopName       string
	}{
		// Valid protocol.
		{dn, 100, 200, NewPixelFormat(), "foo"},
		// Invalid protocol (missing fields).
		{eof: none},
		{eof: fbw, fbWidth: 1},
		{eof: fbh, fbWidth: 2, fbHeight: 1},
		{eof: pf, fbWidth: 3, fbHeight: 2, pixelFormat: NewPixelFormat()},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()
		if tt.eof >= fbw {
			if err := binary.Write(conn.c, binary.BigEndian, tt.fbWidth); err != nil {
				t.Fatal(err)
			}
		}
		if tt.eof >= fbh {
			if err := binary.Write(conn.c, binary.BigEndian, tt.fbHeight); err != nil {
				t.Fatal(err)
			}
		}
		if tt.eof >= pf {
			pfBytes, _ := tt.pixelFormat.Bytes()
			if err := binary.Write(conn.c, binary.BigEndian, pfBytes); err != nil {
				t.Fatal(err)
			}
		}
		if tt.eof >= dn {
			if err := binary.Write(conn.c, binary.BigEndian, uint32(len(tt.desktopName))); err != nil {
				t.Fatal(err)
			}
			if err := binary.Write(conn.c, binary.BigEndian, []byte(tt.desktopName)); err != nil {
				t.Fatal(err)
			}
		}

		// Validate server message handling.
		err := conn.serverInit()
		if tt.eof < dn && err == nil {
			t.Fatalf("serverInit() expected error")
		}
		if tt.eof < dn {
			// The protocol was incomplete; no point in checking values.
			continue
		}
		if err != nil {
			t.Fatalf("serverInit() error %v", err)
		}
		if conn.FrameBufferWidth != tt.fbWidth {
			t.Errorf("serverInit() FrameBufferWidth: got = %v, want = %v", conn.FrameBufferWidth, tt.fbWidth)
		}
		if conn.FrameBufferHeight != tt.fbHeight {
			t.Errorf("serverInit() FrameBufferHeight: got = %v, want = %v", conn.FrameBufferHeight, tt.fbHeight)
		}
		if !equalPixelFormat(conn.PixelFormat, tt.pixelFormat) {
			t.Errorf("serverInit() PixelFormat: got = %v, want = %v", conn.PixelFormat, tt.pixelFormat)
		}
		if conn.DesktopName != tt.desktopName {
			t.Errorf("serverInit() DesktopName: got = %v, want = %v", conn.DesktopName, tt.desktopName)
		}
	}
}
