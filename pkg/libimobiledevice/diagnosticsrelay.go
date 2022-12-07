package libimobiledevice

const (
	DiagnosticsRelayServiceName = "com.apple.mobile.diagnostics_relay"
)

type DiagnosticsRelayBasicRequest struct {
	Request    string  `plist:"Request"`
	Label      string  `plist:"Label"`
	EntryClass *string `plist:"EntryClass,omitempty"`
}

func NewDiagnosticsRelayClient(innerConn InnerConn) *DiagnosticsRelayClient {
	return &DiagnosticsRelayClient{
		newServicePacketClient(innerConn),
	}
}

type DiagnosticsRelayClient struct {
	client *servicePacketClient
}

func (c *DiagnosticsRelayClient) InnerConn() InnerConn {
	return c.client.innerConn
}

func (c *DiagnosticsRelayClient) NewBasicRequest(relayType string, entryClass *string) *DiagnosticsRelayBasicRequest {
	return &DiagnosticsRelayBasicRequest{
		Request:    relayType,
		Label:      BundleID,
		EntryClass: entryClass,
	}
}

func (c *DiagnosticsRelayClient) NewXmlPacket(req interface{}) (Packet, error) {
	return c.client.NewXmlPacket(req)
}

func (c *DiagnosticsRelayClient) SendPacket(pkt Packet) (err error) {
	return c.client.SendPacket(pkt)
}

func (c *DiagnosticsRelayClient) ReceivePacket() (respPkt Packet, err error) {
	return c.client.ReceivePacket()
}
