package lygo_sys

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
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

// ID returns the platform specific machine id of the current host OS.
// Regard the returned id as "confidential" and consider using ProtectedID() instead.
// THANKS TO: github.com/denisbrodbeck/machineid
func ID() (string, error) {
	id, err := machineID()
	if err != nil {
		return "", fmt.Errorf("machineid: %v", err)
	}
	return id, nil
}

// ProtectedID returns a hashed version of the machine ID in a cryptographically secure way,
// using a fixed, application-specific key.
// Internally, this function calculates HMAC-SHA256 of the application ID, keyed by the machine ID.
// THANKS TO: github.com/denisbrodbeck/machineid
func ProtectedID(appID string) (string, error) {
	id, err := ID()
	if err != nil {
		return "", fmt.Errorf("machineid: %v", err)
	}
	return protect(appID, id), nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

// run wraps `exec.Command` with easy access to stdout and stderr.
func run(stdout, stderr io.Writer, cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdin = os.Stdin
	c.Stdout = stdout
	c.Stderr = stderr
	return c.Run()
}

// protect calculates HMAC-SHA256 of the application ID, keyed by the machine ID and returns a hex-encoded string.
func protect(appID, id string) string {
	mac := hmac.New(sha256.New, []byte(id))
	mac.Write([]byte(appID))
	return hex.EncodeToString(mac.Sum(nil))
}

func readFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func trim(s string) string {
	return strings.TrimSpace(strings.Trim(s, "\n"))
}

