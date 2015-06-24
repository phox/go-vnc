package vnc

import (
	"encoding/binary"
	"testing"
)

func TestFramebufferUpdate(t *testing.T) {
	tests := []struct {
		rects []Rectangle
	}{
		{[]Rectangle{{10, 20, 30, 40, NewRawEncoding([]Color{})}}},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		// Send the message.
		msg := NewFramebufferUpdateMessage(tt.rects)
		if err := binary.Write(conn.c, binary.BigEndian, &msg.Msg); err != nil {
			t.Fatal(err)
		}
		for i := range msg.Pad {
			if err := binary.Write(conn.c, binary.BigEndian, &msg.Pad[i]); err != nil {
				t.Fatal(err)
			}
		}
		if err := binary.Write(conn.c, binary.BigEndian, &msg.NumRect); err != nil {
			t.Fatal(err)
		}
		for _, r := range msg.Rects {
			if err := binary.Write(conn.c, binary.BigEndian, &r.X); err != nil {
				t.Fatal(err)
			}
			if err := binary.Write(conn.c, binary.BigEndian, &r.Y); err != nil {
				t.Fatal(err)
			}
			if err := binary.Write(conn.c, binary.BigEndian, &r.Width); err != nil {
				t.Fatal(err)
			}
			if err := binary.Write(conn.c, binary.BigEndian, &r.Height); err != nil {
				t.Fatal(err)
			}
			if err := binary.Write(conn.c, binary.BigEndian, r.Enc.Type()); err != nil {
				t.Fatal(err)
			}
		}

		// Validate message handling.
		var messageType uint8
		if err := binary.Read(conn.c, binary.BigEndian, &messageType); err != nil {
			t.Fatal(err)
		}
		fu := NewFramebufferUpdateMessage([]Rectangle{})
		parsedFU, err := fu.Read(conn, conn.c)
		if err != nil {
			t.Fatalf("FramebufferUpdate() unexpected error %v", err)
		}
		rects := parsedFU.(*FramebufferUpdateMessage).Rects
		if got, want := len(rects), len(tt.rects); got != want {
			t.Errorf("FramebufferUpdate() number of rectangles doesn't match; got = %v, want = %v", got, want)
		}
		for i, r := range tt.rects {
			if got, want := rects[i].X, r.X; got != want {
				t.Errorf("FramebufferUpdate() rect #%v invalid X; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Y, r.Y; got != want {
				t.Errorf("FramebufferUpdate() rect #%v invalid Y; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Width, r.Width; got != want {
				t.Errorf("FramebufferUpdate() rect #%v invalid Width; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Height, r.Height; got != want {
				t.Errorf("FramebufferUpdate() rect #%v invalid Height; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Enc.Type(), r.Enc.Type(); got != want {
				t.Errorf("FramebufferUpdate() rect #%v invalid Enc; got = %v, want = %v", i, got, want)
			}
		}
	}
}

func TestSetColorMapEntries(t *testing.T) {}

func TestBell(t *testing.T) {}

func TestServerCutText(t *testing.T) {}
