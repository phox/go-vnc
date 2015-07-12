package vnc

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func newMockServer(t *testing.T, version string) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("error listening: %s", err)
	}

	go func() {
		defer ln.Close()
		c, err := ln.Accept()
		if err != nil {
			t.Fatalf("error accepting conn: %s", err)
		}
		defer c.Close()

		_, err = c.Write([]byte(fmt.Sprintf("RFB %s\n", version)))
		if err != nil {
			t.Fatal("failed writing version")
		}
	}()

	return ln.Addr().String()
}

func TestLowMajorVersion(t *testing.T) {
	nc, err := net.Dial("tcp", newMockServer(t, "002.009"))
	if err != nil {
		t.Fatalf("error connecting to mock server: %s", err)
	}

	_, err = Connect(context.Background(), nc, &ClientConfig{})
	if err == nil {
		t.Fatal("error expected")
	}
	if err != nil {
		if verr, ok := err.(*VNCError); !ok {
			t.Errorf("Client() unexpected %v error: %v", reflect.TypeOf(err), verr)
		}
	}
}

func TestLowMinorVersion(t *testing.T) {
	nc, err := net.Dial("tcp", newMockServer(t, "003.002"))
	if err != nil {
		t.Fatalf("error connecting to mock server: %s", err)
	}

	_, err = Connect(context.Background(), nc, &ClientConfig{})
	if err == nil {
		t.Fatal("error expected")
	}
	if err != nil {
		if verr, ok := err.(*VNCError); !ok {
			t.Errorf("Client() unexpected %v error: %v", reflect.TypeOf(err), verr)
		}
	}
}

func TestClientConn(t *testing.T) {
	conn := &ClientConn{}

	if got, want := conn.DesktopName(), ""; got != want {
		t.Errorf("DesktopName() failed; got = %v, want = %v", got, want)
	}
	if got, want := conn.FramebufferHeight(), uint16(0); got != want {
		t.Errorf("FramebufferHeight() failed; got = %v, want = %v", got, want)
	}
	if got, want := conn.FramebufferWidth(), uint16(0); got != want {
		t.Errorf("FramebufferWidth() failed; got = %v, want = %v", got, want)
	}
}
