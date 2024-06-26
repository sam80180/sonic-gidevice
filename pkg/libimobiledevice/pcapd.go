package libimobiledevice

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/lunixbochs/struc"
)

const PcapdServiceName = "com.apple.pcapd"

func NewPcapdClient(innerConn InnerConn) *PcapdClient {
	return &PcapdClient{
		client: newServicePacketClient(innerConn),
	}
}

type PcapdClient struct {
	filter func(*IOSPacketHeader) bool
	client *servicePacketClient
}

func (c *PcapdClient) ReceivePacket() (respPkt Packet, err error) {
	var bufLen []byte
	if bufLen, err = c.client.innerConn.Read(4); err != nil {
		return nil, fmt.Errorf("lockdown(Pcapd) receive: %w", err)
	}
	lenPkg := binary.BigEndian.Uint32(bufLen)

	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(bufLen)

	var buf []byte
	if buf, err = c.client.innerConn.Read(int(lenPkg)); err != nil {
		return nil, fmt.Errorf("lockdown(Pcapd) receive: %w", err)
	}
	buffer.Write(buf)

	if respPkt, err = new(servicePacket).Unpack(buffer); err != nil {
		return nil, fmt.Errorf("lockdown(Pcapd) receive: %w", err)
	}

	debugLog(fmt.Sprintf("<-- %s\n", respPkt))

	return
}

type IOSPacketHeader struct {
	HdrSize        uint32  `struc:"uint32,big"`
	Version        uint8   `struc:"uint8,big"`
	PacketSize     uint32  `struc:"uint32,big"`
	Type           uint8   `struc:"uint8,big"`
	Unit           uint16  `struc:"uint16,big"`
	IO             uint8   `struc:"uint8,big"`
	ProtocolFamily uint32  `struc:"uint32,big"`
	FramePreLength uint32  `struc:"uint32,big"`
	FramePstLength uint32  `struc:"uint32,big"`
	IFName         string  `struc:"[16]byte"`
	Pid            int32   `struc:"int32,little"`
	ProcName       string  `struc:"[17]byte"`
	Unknown        uint32  `struc:"uint32,little"`
	Pid2           int32   `struc:"int32,little"`
	ProcName2      string  `struc:"[17]byte"`
	Unknown2       [8]byte `struc:"[8]byte"`
}

var (
	PacketHeaderSize = uint32(95)
)

func (c *PcapdClient) GetPacket(buf []byte) ([]byte, error) {
	iph := IOSPacketHeader{}
	preader := bytes.NewReader(buf)
	_ = struc.Unpack(preader, &iph)

	if c.filter != nil {
		if !c.filter(&iph) {
			return nil, nil
		}
	}

	// This code is from go-ios: https://github.com/danielpaulus/go-ios/blob/main/ios/pcap/pcap.go#164
	if iph.HdrSize > PacketHeaderSize {
		buf := make([]byte, iph.HdrSize-PacketHeaderSize)
		_, err := io.ReadFull(preader, buf)
		if err != nil {
			return []byte{}, err
		}
	}

	packet1, err := ioutil.ReadAll(preader)
	if err != nil {
		return packet1, err
	}
	if iph.FramePreLength == 0 {
		ext := []byte{0xbe, 0xfe, 0xbe, 0xfe, 0xbe, 0xfe, 0xbe, 0xfe, 0xbe, 0xfe, 0xbe, 0xfe, 0x08, 0x00}
		return append(ext, packet1...), nil
	}
	return packet1, nil
}

type PcaprecHdrS struct {
	TsSec   int `struc:"uint32,little"` /* timestamp seconds */
	TsUsec  int `struc:"uint32,little"` /* timestamp microseconds */
	InclLen int `struc:"uint32,little"` /* number of octets of packet saved in file */
	OrigLen int `struc:"uint32,little"` /* actual length of packet */
}

func (c *PcapdClient) CreatePacket(packet []byte) ([]byte, error) {
	now := time.Now()
	phs := &PcaprecHdrS{
		int(now.Unix()),
		int(now.UnixNano()/1e3 - now.Unix()*1e6),
		len(packet),
		len(packet),
	}
	var buf bytes.Buffer
	err := struc.Pack(&buf, phs)
	if err != nil {
		return nil, err
	}

	buf.Write(packet)
	return buf.Bytes(), nil
}

func (c *PcapdClient) Close() {
	c.client.innerConn.Close()
}
