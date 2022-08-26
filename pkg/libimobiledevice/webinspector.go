package libimobiledevice

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"howett.net/plist"
	"strings"
	"time"
)

type WIRMessageKey string

const (
	WebInspectorServiceName               = "com.apple.webinspector"
	WIRFinalMessageKey      WIRMessageKey = "WIRFinalMessageKey"
	WIRPartialMessageKey    WIRMessageKey = "WIRPartialMessageKey"
)

func NewWebInspectorClient(innerConn InnerConn) *WebInspectorClient {
	var web = &WebInspectorClient{
		client:      newServicePacketClient(innerConn),
		MaxPlistLen: 7586,
	}
	web.client.innerConn.Timeout(time.Second * 1)
	return web
}

type WebInspectorClient struct {
	client            *servicePacketClient
	msgBuffer         []string
	MaxPlistLen       int  // 7586
	partialsSupported bool // True for simulators, real devices with iOS version less than 11
}

func (w *WebInspectorClient) SetPartialsSupported(isCompleteSupported bool) {
	w.partialsSupported = isCompleteSupported
}

func (w *WebInspectorClient) NewXmlPacket(req interface{}) (Packet, error) {
	return w.client.NewXmlPacket(req)
}

func (w *WebInspectorClient) SendPacket(pkt Packet) (err error) {
	if !w.partialsSupported {
		// debugLog(fmt.Sprintf("--> Length: %d, Version: %d, Type: %d, Tag: %d\n%s\n", pkt.Length(), pkt.Version(), pkt.Type(), pkt.Tag(), pkt.Body()))
		//debugLog(fmt.Sprintf("--> %s\n", pkt))
		return w.client.SendPacket(pkt)
	} else {
		var raw []byte
		if raw, err = pkt.Pack(); err != nil {
			return fmt.Errorf("send packet: %w", err)
		}
		var isPartial bool
		var partialMsgStr string
		var dataLen = len(raw)
		for i := 0; i < dataLen; i += w.MaxPlistLen {
			var tempSendData = make(map[WIRMessageKey]interface{})
			isPartial = dataLen-i > w.MaxPlistLen
			if isPartial {
				var partialMsgBytes = raw[i : i+w.MaxPlistLen]
				partialMsgStr = string(partialMsgBytes)
				tempSendData[WIRPartialMessageKey] = partialMsgStr
			} else {
				var partialMsgBytes = raw[i:dataLen]
				partialMsgStr = string(partialMsgBytes)
				tempSendData[WIRFinalMessageKey] = partialMsgStr
			}
			temPkg, err := w.NewXmlPacket(tempSendData)
			if err != nil {
				return fmt.Errorf("send packet: %w", err)
			}
			if err = w.client.SendPacket(temPkg); err != nil {
				return fmt.Errorf("send packet: %w", err)
			}
		}
		return err
	}
}

func (w *WebInspectorClient) parseDataAsPlist(respPkt Packet) (response interface{}, err error) {
	var reply interface{}
	if err = respPkt.Unmarshal(&reply); err != nil {
		return nil, fmt.Errorf("receive packet: %w", err)
	}
	return reply, err
}

func (w *WebInspectorClient) SendWebkitMsg(req interface{}) (err error) {
	plistPkt, err := w.NewXmlPacket(req)
	if err != nil {
		return err
	}
	return w.SendPacket(plistPkt)
}

func (w *WebInspectorClient) ReceiveWebkitMsg() (response interface{}, err error) {
	// todo check and debug
	respPkt, err := w.ReceivePacket()
	if err != nil {
		return nil, err
	}
	var reply interface{}
	reply, err = w.parseDataAsPlist(respPkt)
	if err != nil {
		return nil, err
	}
	//debugLog(fmt.Sprintf("<-- %s\n", reply))

	if !w.partialsSupported {
		return reply, nil
	} else {
		var toMapData = reply.(map[WIRMessageKey]interface{})
		var value = toMapData[WIRFinalMessageKey].(string)
		var resultData string
		if value != "" {
			if w.msgBuffer != nil && len(w.msgBuffer) > 0 {
				resultData = strings.Join(w.msgBuffer, "") + value
				w.msgBuffer = w.msgBuffer[:0]
			}
			resultPack, err := w.NewXmlPacket(resultData)
			if err != nil {
				return nil, err
			}
			return w.parseDataAsPlist(resultPack)
		} else {
			value = toMapData[WIRPartialMessageKey].(string)
			w.msgBuffer = append(w.msgBuffer, value)
			return w.ReceiveWebkitMsg()
		}
	}
}

func (w *WebInspectorClient) ReceivePacket() (respPkt Packet, err error) {
	var bufLen []byte
	if bufLen, err = w.client.innerConn.Read(4); err != nil {
		return nil, fmt.Errorf("receive packet: %w", err)
	}
	lenPkg := binary.BigEndian.Uint32(bufLen)

	buffer := bytes.NewBuffer([]byte{})
	buffer.Write(bufLen)

	var buf []byte
	if buf, err = w.client.innerConn.Read(int(lenPkg)); err != nil {
		return nil, fmt.Errorf("receive packet: %w", err)
	}
	buffer.Write(buf)

	if respPkt, err = new(servicePacket).Unpack(buffer); err != nil {
		return nil, fmt.Errorf("receive packet: %w", err)
	}

	debugLog(fmt.Sprintf("<-- %s\n", w.debugRecvPkg(respPkt)))
	return
}

func (w *WebInspectorClient) debugRecvPkg(p Packet) string {
	var reply interface{}
	if err := p.Unmarshal(&reply); err != nil {
		return fmt.Sprintf("recv packet is fail:" + p.String())
	}
	var buf, _ = p.Pack()
	dataString, _ := plist.Marshal(reply, plist.XMLFormat)
	return fmt.Sprintf("Length: %d\n%s", len(buf), dataString)
}
