package vnc

import "testing"

func TestFramebufferUpdate(t *testing.T) {
	tests := []struct {
		rects []RectangleMessage
	}{
		{[]RectangleMessage{{10, 20, 30, 40, Raw}}},
	}

	mockConn := &MockConn{}
	conn := &ClientConn{
		c:      mockConn,
		config: &ClientConfig{},
	}

	for _, tt := range tests {
		mockConn.Reset()

		// Send the message.
		msg := &FramebufferUpdateMessage{
			Msg:     FramebufferUpdateMsg,
			NumRect: uint16(len(tt.rects)),
		}
		if err := conn.send(msg); err != nil {
			t.Error(err)
			continue
		}
		if err := conn.send(tt.rects); err != nil {
			t.Error(err)
			continue
		}

		// Validate message handling.
		var messageType uint8
		if err := conn.receive(&messageType); err != nil {
			t.Error(err)
			continue
		}
		fu := &FramebufferUpdate{}
		parsedFU, err := fu.Read(conn, conn.c)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		rects := parsedFU.(*FramebufferUpdate).Rects

		// Validate the message.
		if got, want := len(rects), len(tt.rects); got != want {
			t.Errorf("incorrect number-of-rectangles; got = %v, want = %v", got, want)
			continue
		}
		for i, r := range tt.rects {
			if got, want := rects[i].X, r.X; got != want {
				t.Errorf("rect[%v] incorrect x-position; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Y, r.Y; got != want {
				t.Errorf("rect[%v] incorrect y-position; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Width, r.Width; got != want {
				t.Errorf("rect[%v] incorrect width; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Height, r.Height; got != want {
				t.Errorf("rect[%v] incorrect height; got = %v, want = %v", i, got, want)
			}
			if got, want := rects[i].Enc.Type(), r.Encoding; got != want {
				t.Errorf("rect[%v] incorrect encoding-type; got = %v, want = %v", i, got, want)
			}
		}
	}
}

func TestSetColorMapEntries(t *testing.T) {}

func TestBell(t *testing.T) {}

func TestServerCutText(t *testing.T) {}
