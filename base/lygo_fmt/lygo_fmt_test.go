package lygo_fmt

import (
	"fmt"
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	now := time.Now()

	patterns := []string{
		"yyyy-MM-dd",
		"yyyy-MM-dd HH:mm:ss",
		"yyyy-MM-dd HH:mm:ssZ",
	}

	for _, pattern:=range patterns{
		f := FormatDate(now, pattern)
		fmt.Println(f)
		d, err := ParseDate(f, pattern)
		if nil!=err{
			t.Error(err)
			t.FailNow()
		}
		fmt.Println(d)
	}

}
