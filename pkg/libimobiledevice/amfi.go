package libimobiledevice

import (
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

const (
	AmfiServiceName = "com.apple.amfi.lockdown"

	DEV_MODE_REVEAL = 0
	DEV_MODE_ARM    = 1
	DEV_MODE_ENABLE = 2
)

type AmfiClient struct {
	client *servicePacketClient
}

func NewAmfiClient(innerConn InnerConn) *AmfiClient {
	return &AmfiClient{
		newServicePacketClient(innerConn),
	}
}

func (c *AmfiClient) SendAction(action int) (int, error) {
	data := map[string]int{"action": action}
	pktOut, _ := c.client.NewXmlPacket(data)
	if errSend := c.client.SendPacket(pktOut); errSend != nil {
		return http.StatusInternalServerError, errSend
	}
	if pktIn, errRecv := c.client.ReceivePacket(); errRecv != nil {
		return http.StatusInternalServerError, errRecv
	} else {
		result := map[string]interface{}{}
		if errDec := pktIn.Unmarshal(&result); errDec != nil {
			return http.StatusInternalServerError, errDec
		}
		if oErr, bHasErr := result["Error"]; bHasErr {
			switch strErr := oErr.(type) {
			case string:
				if strings.Contains(strErr, "passcode") {
					return http.StatusLocked, xerrors.New(strErr)
				}
				return http.StatusInternalServerError, xerrors.New(strErr)
			}
		}
		if oSuccess, bHasSuccess := result["success"]; bHasSuccess {
			switch bSuccess := oSuccess.(type) {
			case bool:
				if bSuccess {
					return http.StatusOK, nil
				}
				return http.StatusConflict, nil
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				if bSuccess != 0 {
					return http.StatusOK, nil
				}
				return http.StatusConflict, nil
			case string:
				if bSuccess == "true" || bSuccess == "1" {
					return http.StatusOK, nil
				}
				return http.StatusConflict, nil
			}
		}
		return http.StatusUnprocessableEntity, xerrors.Errorf("Unknown response: %+v", result)
	}
}
