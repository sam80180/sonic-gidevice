package giDevice

import (
	"github.com/SonicCloudOrg/sonic-gidevice/pkg/libimobiledevice"
)

func newDiagnosticsRelay(client *libimobiledevice.DiagnosticsRelayClient) *diagnostics {
	return &diagnostics{
		client: client,
	}
}

type diagnostics struct {
	client *libimobiledevice.DiagnosticsRelayClient
}

func (d *diagnostics) Reboot() (err error) {
	var pkt libimobiledevice.Packet
	if pkt, err = d.client.NewXmlPacket(
		d.client.NewBasicRequest("Restart"),
	); err != nil {
		return
	}
	if err = d.client.SendPacket(pkt); err != nil {
		return err
	}
	return
}

func (d *diagnostics) Shutdown() (err error) {
	var pkt libimobiledevice.Packet
	if pkt, err = d.client.NewXmlPacket(
		d.client.NewBasicRequest("Shutdown"),
	); err != nil {
		return
	}
	if err = d.client.SendPacket(pkt); err != nil {
		return err
	}
	return
}

func (d *diagnostics) PowerSource() (powerInfo interface{}, err error) {
	var data map[string]string
	data["Request"] = "IORegistry"
	data["EntryClass"] = "IOPMPowerSource"
	data["Label"] = "org.cloud.sonic.gidevice"
	var pkt libimobiledevice.Packet
	if pkt, err = d.client.NewXmlPacket(
		d.client.NewBasicRequest("Shutdown"),
	); err != nil {
		return nil, err
	}
	if err = d.client.SendPacket(pkt); err != nil {
		return nil, err
	}
	if pkt, err = d.client.ReceivePacket(); err != nil {
		return nil, err
	}
	if err = pkt.Unmarshal(powerInfo); err != nil {
		return nil, err
	}
	return
}
