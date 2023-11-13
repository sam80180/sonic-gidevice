package giDevice

import (
	"github.com/SonicCloudOrg/sonic-gidevice/pkg/libimobiledevice"
)

type amfi struct {
	client *libimobiledevice.AmfiClient
}

var _ Amfi = (*amfi)(nil)

func newAmfi(client *libimobiledevice.AmfiClient) *amfi {
	return &amfi{client: client}
}

func (c *amfi) DevModeReveal() (int, error) {
	return c.client.SendAction(libimobiledevice.DEV_MODE_REVEAL)
}

func (c *amfi) DevModeArm() (int, error) {
	return c.client.SendAction(libimobiledevice.DEV_MODE_ARM)
}

func (c *amfi) DevModeEnable() (int, error) {
	return c.client.SendAction(libimobiledevice.DEV_MODE_ENABLE)
}
