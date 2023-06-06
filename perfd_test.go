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
	"testing"
	"time"
)

func TestPerfSystemMonitor(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(true),
		WithPerfSystemMem(true),
		WithPerfSystemDisk(true),
		WithPerfSystemNetwork(true),
		WithPerfOutputInterval(2000),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 20))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfSystemCpu(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(true),
		WithPerfOutputInterval(1000),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfSystemMem(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemMem(true),
		WithPerfOutputInterval(1000),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfNotSystemPerfData(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemMem(false),
		WithPerfSystemCPU(false),
		WithPerfOutputInterval(1000),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfProcessMonitor(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(true),
		WithPerfProcessAttributes("cpuUsage", "memAnon"),
		WithPerfOutputInterval(1000),
		WithPerfPID(100),
		WithPerfBundleID("com.apple.mobilesafari"), // higher priority than pid
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfGPU(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(false),
		WithPerfSystemMem(false),
		WithPerfGPU(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfFPS(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(false),
		WithPerfSystemMem(false),
		WithPerfFPS(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfNetwork(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(false),
		WithPerfSystemMem(false),
		WithPerfNetwork(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}

func TestPerfAll(t *testing.T) {
	setupLockdownSrv(t)

	data, err := dev.PerfStart(
		WithPerfSystemCPU(true),
		WithPerfSystemMem(true),
		WithPerfSystemDisk(true),
		WithPerfSystemNetwork(true),
		WithPerfNetwork(true),
		WithPerfFPS(true),
		WithPerfGPU(true),
		WithPerfProcessAttributes("cpuUsage", "memAnon"),
		WithPerfBundleID("com.apple.mobilesafari"),
	)
	if err != nil {
		t.Fatal(err)
	}

	timer := time.NewTimer(time.Duration(time.Second * 10))
	for {
		select {
		case <-timer.C:
			dev.PerfStop()
			return
		case d := <-data:
			fmt.Println(string(d))
		}
	}
}
