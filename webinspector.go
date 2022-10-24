package giDevice

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-gidevice/pkg/libimobiledevice"
)

var _ WebInspector = (*webInspectorService)(nil)

func newWebInspectorService(client *libimobiledevice.WebInspectorClient) *webInspectorService {
	return &webInspectorService{
		client: client,
	}
}

type webInspectorService struct {
	client *libimobiledevice.WebInspectorClient
}

func (w *webInspectorService) SetPartialsSupported(isCompleteSupported bool) {
	w.client.SetPartialsSupported(isCompleteSupported)
}

func (w *webInspectorService) SetPartialsMaxLength(maxLent int) {
	w.client.MaxPlistLen = maxLent
}

func (w *webInspectorService) SendWebkitMsg(selector string, args interface{}) error {
	if selector == "" {
		return fmt.Errorf("selector cannot be empty")
	}
	var request = make(map[string]interface{})
	request["__selector"] = selector
	request["__argument"] = args

	err := w.client.SendWebkitMsg(request)
	if err != nil {
		return err
	}
	return nil
}

func (w *webInspectorService) ReceiveWebkitMsg() (response interface{}, err error) {
	return w.client.ReceiveWebkitMsg()
}
