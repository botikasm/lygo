package lygo_exec

import (
	"errors"
	"os/exec"
)

// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Run(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Open(args ...string) error {
	c := exec.Command(OPEN_FILE_COMMAND, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return err
	}
	if len(out) > 0 {
		// may be an error
		return errors.New(string(out))
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
