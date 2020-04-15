// +build darwin

package lygo_sys

import (
	"os/exec"
)

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func shutdown(adminPsw string) error{
	// echo <password> | sudo -S shutdown -h now
	if err := exec.Command("/bin/sh", "-c", "echo " +  adminPsw + " | sudo -S shutdown -h now").Run(); err != nil {
		return err
	}
	return nil
}

