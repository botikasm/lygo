package lygo_mu

import "github.com/botikasm/lygo/base/lygo_conv"

var (
	// BYTES
	Kb = 1024
	Mb = Kb * 1024
	Gb = Mb * 1024
	Tb = Gb * 1024
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

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