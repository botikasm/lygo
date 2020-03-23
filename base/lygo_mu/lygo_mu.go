package lygo_mu

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
)

// ---------------------------------------------------------------------------------------------------------------------
// Measure Unity utility.
// Here are some utility methods to easily convert or format measure unit.
// ---------------------------------------------------------------------------------------------------------------------

var (
	// BYTES
	Kb = int64(1024)
	Mb = Kb * 1024
	Gb = Mb * 1024
	Tb = Gb * 1024
	Pb = Tb * 1024
	Eb = Pb * 1024
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func FmtBytes(b int64) string {
	const unit = 1024
	if b < Kb {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func ToKiloBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Kb)
}

func ToMegaBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Mb)
}

func ToGigaBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Gb)
}

func ToTeraBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Tb)
}

func ToPetaBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Pb)
}

func ToEsaBytes(val int64) float64 {
	return lygo_conv.ToFloat64(val) / lygo_conv.ToFloat64(Eb)
}