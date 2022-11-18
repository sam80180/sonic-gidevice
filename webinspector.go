/*
 *  sonic-gidevice  Connect to your iOS Devices.
 *  Copyright (C) 2022 SonicCloudOrg
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
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
