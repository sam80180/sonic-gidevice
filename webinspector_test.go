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
