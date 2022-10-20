package giDevice

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"testing"
)

type NetworkInfo struct {
	Mac  string
	IPv4 string
	IPv6 string
}

func TestPacd(t *testing.T) {
	setupLockdownSrv(t)
	GetNetworkIP()
}

func GetNetworkIP() error {
	mac, err := dev.GetValue("", "WiFiAddress")
	if err != nil {
		fmt.Println("not key")
	}
	macStr, _ := mac.(string)
	info := NetworkInfo{}
	info.Mac = macStr
	resultBytes, err := dev.Pcap()
	if err != nil {
		return err
	}
	for {
		select {
		case data, ok := <-resultBytes:
			if ok {
				err = findIP(data, &info)
				if err != nil {
					return err
				}
				if info.Mac != "" && info.IPv6 != "" && info.IPv4 != "" {
					return nil
				}
			}
			log.Println(&info)
		}
	}
}

func findIP(p []byte, info *NetworkInfo) error {
	packet := gopacket.NewPacket(p, layers.LayerTypeEthernet, gopacket.Default)
	// Get the TCP layer from this packet
	if tcpLayer := packet.Layer(layers.LayerTypeEthernet); tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.Ethernet)
		if tcp.SrcMAC.String() == info.Mac {
			if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
				ipv4, ok := ipv4Layer.(*layers.IPv4)
				if ok {
					info.IPv4 = ipv4.SrcIP.String()
					log.Printf("ip4 found:%s", info.IPv4)
				}
			}
			if ipv6Layer := packet.Layer(layers.LayerTypeIPv6); ipv6Layer != nil {
				ipv6, ok := ipv6Layer.(*layers.IPv6)
				if ok {
					info.IPv6 = ipv6.SrcIP.String()
					log.Printf("ip6 found:%s", info.IPv6)
				}
			}
		}
	}
	return nil
}
