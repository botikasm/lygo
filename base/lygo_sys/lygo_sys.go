package lygo_sys

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_json"
	"runtime"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type InfoObject struct {
	GoOS     string `json:"goos"`
	Kernel   string `json:"kernel"`
	Core     string `json:"core"`
	Platform string `json:"platform"`
	OS       string `json:"os"`
	Hostname string `json:"hostname"`
	CPUs     int    `json:"cpus"`
}

func (instance *InfoObject) VarDump() {
	fmt.Println("GoOS:", instance.GoOS)
	fmt.Println("Kernel:", instance.Kernel)
	fmt.Println("Core:", instance.Core)
	fmt.Println("Platform:", instance.Platform)
	fmt.Println("OS:", instance.OS)
	fmt.Println("Hostname:", instance.Hostname)
	fmt.Println("CPUs:", instance.CPUs)
}

func (instance *InfoObject) ToString() string {
	return fmt.Sprintf("GoOS:%v,Kernel:%v,Core:%v,Platform:%v,OS:%v,Hostname:%v,CPUs:%v", instance.GoOS, instance.Kernel, instance.Core, instance.Platform, instance.OS, instance.Hostname, instance.CPUs)
}

func (instance *InfoObject) ToJsonString() string {
	return lygo_json.Stringify(instance)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func GetInfo() *InfoObject {
	return getInfo()
}

func GetOS() string {
	return runtime.GOOS
}

func IsMac() bool {
	return "darwin" == GetOS()
}

func IsLinux() bool {
	return "linux" == GetOS()
}

func IsWindows() bool {
	return "windows" == GetOS()
}

func GetOSVersion() string {
	return GetInfo().Core
}

// shutdown the machine
func Shutdown(a ...string) error {
	adminPsw := ""
	if len(a) == 1 {
		adminPsw = a[0]
	}
	return shutdown(adminPsw)
}
