// Implementation of RFC 6143 §7.3 Initialization Messages.

package vnc

import "log"

// clientInit implements §7.3.1 ClientInit.
func (c *ClientConn) clientInit() error {
	if c.debug {
		log.Print("clientInit()")
	}

	sharedFlag := uint8(0)
	if !c.config.Exclusive {
		sharedFlag = 1
	}
	if c.debug {
		log.Printf("sharedFlag: %v", sharedFlag)
	}
	if err := c.send(sharedFlag); err != nil {
		return err
	}

	return nil
}

// serverInit implements §7.3.2 ServerInit.
func (c *ClientConn) serverInit() error {
	if c.debug {
		log.Print("serverInit()")
	}

	var width, height uint16
	if err := c.receive(&width); err != nil {
		return err
	}
	if err := c.receive(&height); err != nil {
		return err
	}
	c.setFramebufferWidth(width)
	c.setFramebufferHeight(height)

	if err := c.pixelFormat.Read(c.c); err != nil {
		return err
	}

	var nameLength uint32
	if err := c.receive(&nameLength); err != nil {
		return err
	}
	nameBytes := make([]uint8, nameLength)
	if err := c.receive(&nameBytes); err != nil {
		return err
	}
	c.setDesktopName(string(nameBytes))

	return nil
}
