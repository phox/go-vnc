// Implementation of RFC 6143 §7.3 Initialization Messages.

package vnc

import (
	"github.com/golang/glog"
	"github.com/kward/go-vnc/logging"
)

// clientInit implements §7.3.1 ClientInit.
func (c *ClientConn) clientInit() error {
	if glog.V(logging.FnDeclLevel) {
		glog.Info(logging.FnName())
	}

	sharedFlag := uint8(0)
	if !c.config.Exclusive {
		sharedFlag = RFBTrue
	}
	if glog.V(logging.ResultLevel) {
		glog.Infof("sharedFlag: %d", sharedFlag)
	}
	if err := c.send(sharedFlag); err != nil {
		return err
	}

	return nil
}

// serverInit implements §7.3.2 ServerInit.
func (c *ClientConn) serverInit() error {
	if glog.V(logging.FnDeclLevel) {
		glog.Info(logging.FnName())
	}

	var width, height uint16
	if err := c.receive(&width); err != nil {
		return err
	}
	c.setFramebufferWidth(width)
	if err := c.receive(&height); err != nil {
		return err
	}
	if glog.V(logging.ResultLevel) {
		glog.Infof("width: %d height: %d", width, height)
	}
	c.setFramebufferHeight(height)

	if err := c.pixelFormat.Read(c.c); err != nil {
		return err
	}
	if glog.V(logging.ResultLevel) {
		glog.Infof("pixelFormat: %v", c.pixelFormat)
	}

	var nameLength uint32
	if err := c.receive(&nameLength); err != nil {
		return err
	}
	if glog.V(logging.ResultLevel) {
		glog.Infof("desktopName length: %d", nameLength)
	}

	name := make([]uint8, nameLength)
	if err := c.receive(&name); err != nil {
		return err
	}
	if glog.V(logging.ResultLevel) {
		glog.Infof("desktopName: %s", name)
	}
	c.setDesktopName(string(name))

	return nil
}
