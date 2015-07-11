/*
initialization.go implements RFC 6143 §7.3 Initialization Messages.
See http://tools.ietf.org/html/rfc6143#section-7.3 for more info.
*/
package vnc

// clientInit implements §7.3.1 ClientInit.
func (c *ClientConn) clientInit() error {
	var sharedFlag uint8

	if !c.config.Exclusive {
		sharedFlag = 1
	}
	if err := c.send(sharedFlag); err != nil {
		return err
	}

	return nil
}

// serverInit implements §7.3.2 ServerInit.
func (c *ClientConn) serverInit() error {
	if err := c.receive(&c.fbWidth); err != nil {
		return err
	}
	if err := c.receive(&c.fbHeight); err != nil {
		return err
	}
	if err := c.pixelFormat.Write(c.c); err != nil {
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
	c.desktopName = string(nameBytes)

	return nil
}
