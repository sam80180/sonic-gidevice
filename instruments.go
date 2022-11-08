package giDevice

import (
	"encoding/json"
	"fmt"

	"github.com/SonicCloudOrg/sonic-gidevice/pkg/libimobiledevice"
)

// instruments services
const (
	instrumentsServiceDeviceInfo              = "com.apple.instruments.server.services.deviceinfo"
	instrumentsServiceProcessControl          = "com.apple.instruments.server.services.processcontrol"
	instrumentsServiceDeviceApplictionListing = "com.apple.instruments.server.services.device.applictionListing"
	instrumentsServiceGraphicsOpengl          = "com.apple.instruments.server.services.graphics.opengl"     // 获取 GPU/FPS
	instrumentsServiceSysmontap               = "com.apple.instruments.server.services.sysmontap"           // 获取 CPU/Mem/Disk/Network 性能数据
	instrumentsServiceNetworking              = "com.apple.instruments.server.services.networking"          // 获取所有网络详情数据
	instrumentsServiceMobileNotifications     = "com.apple.instruments.server.services.mobilenotifications" // 监控应用状态
)

const (
	instrumentsServiceXcodeNetworkStatistics = "com.apple.xcode.debug-gauge-data-providers.NetworkStatistics" // 获取单进程网络数据
	instrumentsServiceXcodeEnergyStatistics  = "com.apple.xcode.debug-gauge-data-providers.Energy"            // 获取功耗数据
)

var _ Instruments = (*instruments)(nil)

func newInstruments(client *libimobiledevice.InstrumentsClient) *instruments {
	return &instruments{
		client: client,
	}
}

type instruments struct {
	client *libimobiledevice.InstrumentsClient
}

func (i *instruments) notifyOfPublishedCapabilities() (err error) {
	_, err = i.client.NotifyOfPublishedCapabilities()
	return
}

func (i *instruments) requestChannel(channel string) (id uint32, err error) {
	return i.client.RequestChannel(channel)
}

func (i *instruments) AppLaunch(bundleID string, opts ...AppLaunchOption) (pid int, err error) {
	opt := new(appLaunchOption)
	opt.appPath = ""
	opt.options = map[string]interface{}{
		"StartSuspendedKey": uint64(0),
		"KillExisting":      uint64(0),
	}
	if len(opts) != 0 {
		for _, optFunc := range opts {
			optFunc(opt)
		}
	}

	var id uint32
	if id, err = i.requestChannel(instrumentsServiceProcessControl); err != nil {
		return 0, err
	}

	args := libimobiledevice.NewAuxBuffer()
	if err = args.AppendObject(opt.appPath); err != nil {
		return 0, err
	}
	if err = args.AppendObject(bundleID); err != nil {
		return 0, err
	}
	if err = args.AppendObject(opt.environment); err != nil {
		return 0, err
	}
	if err = args.AppendObject(opt.arguments); err != nil {
		return 0, err
	}
	if err = args.AppendObject(opt.options); err != nil {
		return 0, err
	}

	var result *libimobiledevice.DTXMessageResult
	selector := "launchSuspendedProcessWithDevicePath:bundleIdentifier:environment:arguments:options:"
	if result, err = i.client.Invoke(selector, args, id, true); err != nil {
		return 0, err
	}

	if nsErr, ok := result.Obj.(libimobiledevice.NSError); ok {
		return 0, fmt.Errorf("%s", nsErr.NSUserInfo.(map[string]interface{})["NSLocalizedDescription"])
	}

	return int(result.Obj.(uint64)), nil
}

func (i *instruments) appProcess(bundleID string) (err error) {
	var id uint32
	if id, err = i.requestChannel(instrumentsServiceProcessControl); err != nil {
		return err
	}

	args := libimobiledevice.NewAuxBuffer()
	if err = args.AppendObject(bundleID); err != nil {
		return err
	}

	selector := "processIdentifierForBundleIdentifier:"
	if _, err = i.client.Invoke(selector, args, id, true); err != nil {
		return err
	}

	return
}

func (i *instruments) startObserving(pid int) (err error) {
	var id uint32
	if id, err = i.requestChannel(instrumentsServiceProcessControl); err != nil {
		return err
	}

	args := libimobiledevice.NewAuxBuffer()
	if err = args.AppendObject(pid); err != nil {
		return err
	}

	var result *libimobiledevice.DTXMessageResult
	selector := "startObservingPid:"
	if result, err = i.client.Invoke(selector, args, id, true); err != nil {
		return err
	}

	if nsErr, ok := result.Obj.(libimobiledevice.NSError); ok {
		return fmt.Errorf("%s", nsErr.NSUserInfo.(map[string]interface{})["NSLocalizedDescription"])
	}
	return
}

func (i *instruments) AppKill(pid int) (err error) {
	var id uint32
	if id, err = i.requestChannel(instrumentsServiceProcessControl); err != nil {
		return err
	}

	args := libimobiledevice.NewAuxBuffer()
	if err = args.AppendObject(pid); err != nil {
		return err
	}

	selector := "killPid:"
	if _, err = i.client.Invoke(selector, args, id, false); err != nil {
		return err
	}

	return
}

func (i *instruments) AppRunningProcesses() (processes []Process, err error) {
	var id uint32
	if id, err = i.requestChannel(instrumentsServiceDeviceInfo); err != nil {
		return nil, err
	}

	selector := "runningProcesses"

	var result *libimobiledevice.DTXMessageResult
	if result, err = i.client.Invoke(selector, libimobiledevice.NewAuxBuffer(), id, true); err != nil {
		return nil, err
	}

	objs := result.Obj.([]interface{})

	processes = make([]Process, 0, len(objs))

	for _, v := range objs {
		m := v.(map[string]interface{})

		var data []byte
		if data, err = json.Marshal(m); err != nil {
			debugLog(fmt.Sprintf("process marshal: %v\n%v\n", err, m))
			err = nil
			continue
		}

		var tp Process
		if err = json.Unmarshal(data, &tp); err != nil {
			debugLog(fmt.Sprintf("process unmarshal: %v\n%v\n", err, m))
			err = nil
			continue
		}

		processes = append(processes, tp)
	}

	return
}

func (i *instruments) AppList(opts ...AppListOption) (apps []Application, err error) {
	opt := new(appListOption)
	opt.updateToken = ""
	opt.appsMatching = make(map[string]interface{})
	if len(opts) != 0 {
		for _, optFunc := range opts {
			optFunc(opt)
		}
	}

	var id uint32
	if id, err = i.requestChannel(instrumentsServiceDeviceApplictionListing); err != nil {
		return nil, err
	}

	args := libimobiledevice.NewAuxBuffer()
	if err = args.AppendObject(opt.appsMatching); err != nil {
		return nil, err
	}
	if err = args.AppendObject(opt.updateToken); err != nil {
		return nil, err
	}

	selector := "installedApplicationsMatching:registerUpdateToken:"

	var result *libimobiledevice.DTXMessageResult
	if result, err = i.client.Invoke(selector, args, id, true); err != nil {
		return nil, err
	}

	objs := result.Obj.([]interface{})

	for _, v := range objs {
		m := v.(map[string]interface{})

		var data []byte
		if data, err = json.Marshal(m); err != nil {
			debugLog(fmt.Sprintf("application marshal: %v\n%v\n", err, m))
			err = nil
			continue
		}

		var app Application
		if err = json.Unmarshal(data, &app); err != nil {
			debugLog(fmt.Sprintf("application unmarshal: %v\n%v\n", err, m))
			err = nil
			continue
		}
		apps = append(apps, app)
	}

	return
}

func (i *instruments) DeviceInfo() (devInfo *DeviceInfo, err error) {
	var id uint32
	if id, err = i.requestChannel(instrumentsServiceDeviceInfo); err != nil {
		return nil, err
	}

	selector := "systemInformation"

	var result *libimobiledevice.DTXMessageResult
	if result, err = i.client.Invoke(selector, libimobiledevice.NewAuxBuffer(), id, true); err != nil {
		return nil, err
	}

	data, err := json.Marshal(result.Obj)
	if err != nil {
		return nil, err
	}
	devInfo = new(DeviceInfo)
	err = json.Unmarshal(data, devInfo)

	return
}

func (i *instruments) registerCallback(obj string, cb func(m libimobiledevice.DTXMessageResult)) {
	i.client.RegisterCallback(obj, cb)
}

func (i *instruments) call(channel, selector string, auxiliaries ...interface{}) (
	result *libimobiledevice.DTXMessageResult, err error) {

	chanID, err := i.requestChannel(channel)
	if err != nil {
		return nil, err
	}

	args := libimobiledevice.NewAuxBuffer()
	for _, aux := range auxiliaries {
		if err = args.AppendObject(aux); err != nil {
			return nil, err
		}
	}

	return i.client.Invoke(selector, args, chanID, true)
}

func (i *instruments) getPidByBundleID(bundleID string) (pid int, err error) {
	apps, err := i.AppList()
	if err != nil {
		fmt.Printf("get app list error: %v\n", err)
		return 0, err
	}

	mapper := make(map[string]interface{})
	for _, app := range apps {
		mapper[app.ExecutableName] = app.CFBundleIdentifier
	}

	processes, err := i.AppRunningProcesses()
	if err != nil {
		fmt.Printf("get running app processes error: %v\n", err)
		return 0, err
	}
	for _, proc := range processes {
		b, ok := mapper[proc.Name]
		if ok && bundleID == b {
			fmt.Printf("get pid %d by bundleId %s\n", proc.Pid, bundleID)
			return proc.Pid, nil
		}
	}

	fmt.Printf("can't find pid by bundleID: %s\n", bundleID)
	return 0, fmt.Errorf("can't find pid by bundleID: %s", bundleID)
}

type Application struct {
	AppExtensionUUIDs         []string `json:"AppExtensionUUIDs,omitempty"`
	BundlePath                string   `json:"BundlePath"`
	CFBundleIdentifier        string   `json:"CFBundleIdentifier"`
	ContainerBundleIdentifier string   `json:"ContainerBundleIdentifier,omitempty"`
	ContainerBundlePath       string   `json:"ContainerBundlePath,omitempty"`
	DisplayName               string   `json:"DisplayName"`
	ExecutableName            string   `json:"ExecutableName,omitempty"`
	Placeholder               bool     `json:"Placeholder,omitempty"`
	PluginIdentifier          string   `json:"PluginIdentifier,omitempty"`
	PluginUUID                string   `json:"PluginUUID,omitempty"`
	Restricted                int      `json:"Restricted"`
	Type                      string   `json:"Type"`
	Version                   string   `json:"Version"`
}

type DeviceInfo struct {
	Description       string `json:"_deviceDescription"`
	DisplayName       string `json:"_deviceDisplayName"`
	Identifier        string `json:"_deviceIdentifier"`
	Version           string `json:"_deviceVersion"`
	ProductType       string `json:"_productType"`
	ProductVersion    string `json:"_productVersion"`
	XRDeviceClassName string `json:"_xrdeviceClassName"`
}

type CPUInfo struct {
	Pid            string  `json:"PID,omitempty"`                   // 进程ID
	CPUCount       int     `json:"SystemCpuCount"`                  // CPU总数
	TimeStamp      int64   `json:"TimeStamp"`                       // 时间戳
	CPUUsage       float64 `json:"ProcessCpuUsage,omitempty"`       // 单个进程的CPU使用率
	SysCpuUsage    float64 `json:"SystemCpuUsage,omitempty"`        // 系统总体CPU占用
	AttrCtxSwitch  int64   `json:"ProcessAttrCtxSwitch,omitempty"`  // 上下文切换数
	AttrIntWakeups int64   `json:"ProcessAttrIntWakeups,omitempty"` // 唤醒数
	Mess           string  `json:"ProcessInfo,omitempty"`           // 提示信息，当PID没输入或者信息错误时提示
}

type FPSInfo struct {
	FPS       int   `json:"FPS"`
	TimeStamp int64 `json:"TimeStamp"` // 时间戳
}

type GPUInfo struct {
	TilerUtilization    int64 `json:"TilerUtilization"`    // 处理顶点的GPU时间占比
	TimeStamp           int64 `json:"TimeStamp"`           // 时间戳
	DeviceUtilization   int64 `json:"DeviceUtilization"`   // 设备利用率
	RendererUtilization int64 `json:"RendererUtilization"` // 渲染器利用率
}

type MEMInfo struct {
	Anon       int64  `json:"AnonMemory"`            // 虚拟内存
	PhysMemory int64  `json:"PhysMemory"`            // 物理内存
	Rss        int64  `json:"Rss"`                   // 总内存
	Vss        int64  `json:"Vss"`                   // 虚拟内存
	TimeStamp  int64  `json:"TimeStamp"`             // 时间戳
	Mess       string `json:"ProcessInfo,omitempty"` // 提示信息，当PID没输入或者信息错误时提示
}

type NetWorkingInfo struct {
	RxBytes   int64 `json:"RxBytes"`
	RxPackets int64 `json:"RxPackets"`
	TxBytes   int64 `json:"TxBytes"`
	TxPackets int64 `json:"TxPackets"`
	TimeStamp int64 `json:"TimeStamp"` // 时间戳
}
