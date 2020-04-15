// +build windows

package lygo_sys

import "os/exec"

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func shutdown(adminPsw string) error {
	// cmd := "shutdown -s -t O"
	if err := exec.Command("cmd", "shutdown", "-s", "-t", "O").Run(); err != nil {
		//fmt.Println("Failed to initiate shutdown:", err)
		return err
	}
	return nil
}
