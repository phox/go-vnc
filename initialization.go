/*
initialization.go implements RFC 6143 §7.3 Initialization Messages.
See http://tools.ietf.org/html/rfc6143#section-7.3 for more info.
*/
package vnc

import (
	"encoding/binary"
)

// clientInit implements §7.3.1 ClientInit.
func (c *ClientConn) clientInit() error {
	var sharedFlag uint8

	if !c.config.Exclusive {
		sharedFlag = 1
	}
	if err := binary.Write(c.c, binary.BigEndian, sharedFlag); err != nil {
		return err
	}

	return nil
}

// serverInit implements §7.3.2 ServerInit.
func (c *ClientConn) serverInit() error {
	if err := binary.Read(c.c, binary.BigEndian, &c.FramebufferWidth); err != nil {
		return err
	}
	if err := binary.Read(c.c, binary.BigEndian, &c.FramebufferHeight); err != nil {
		return err
	}
	if err := c.PixelFormat.Write(c.c); err != nil {
		return err
	}

	var nameLength uint32
	if err := binary.Read(c.c, binary.BigEndian, &nameLength); err != nil {
		return err
	}
	nameBytes := make([]uint8, nameLength)
	if err := binary.Read(c.c, binary.BigEndian, &nameBytes); err != nil {
		return err
	}
	c.desktopName = string(nameBytes)

	return nil
}
