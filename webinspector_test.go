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
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

var webInspector WebInspector

func setupWebInspectorSrv(t *testing.T) {
	setupLockdownSrv(t)

	var err error
	if lockdownSrv, err = dev.lockdownService(); err != nil {
		t.Fatal(err)
	}

	if webInspector, err = lockdownSrv.WebInspectorService(); err != nil {
		t.Fatal(err)
	}
}

func Test_webInspector_Connect(t *testing.T) {
	setupWebInspectorSrv(t)
	//SetDebug(true,true)
	webInspector.SetPartialsSupported(false)
	var connectID = strings.ToUpper("0330e5d7-45e5-4bf7-96bc-c5a3afd4615b")
	var data = map[string]interface{}{
		"WIRConnectionIdentifierKey": connectID,
	}
	webInspector.SendWebkitMsg("_rpc_reportIdentifier:", data)
	start := time.Now().Unix()
	ctx, cancnel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("end")
				return
			default:
				if response, err := webInspector.ReceiveWebkitMsg(); err != nil {
					if !strings.Contains(err.Error(), "timeout") {
						log.Fatal(err)
						return
					}
				} else {
					fmt.Println("==============recv connect==================")
					fmt.Println(response)
				}
				if time.Now().Unix()-start > 4 {
					cancnel()
				}
			}
		}
	}()
	time.Sleep(4 * time.Second)
}
