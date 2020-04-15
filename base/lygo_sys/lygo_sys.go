package lygo_sys

import "runtime"

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

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

// shutdown the machine
func Shutdown(a ...string) error {
	adminPsw := ""
	if len(a) == 1 {
		adminPsw = a[0]
	}
	return shutdown(adminPsw)
}
