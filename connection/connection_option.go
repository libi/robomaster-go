package connection

import "time"

var defaultOption = &Option{
	ScanTimeout: time.Second * 5,
	CtrlTimeOut: time.Second,
}

type Option struct {
	IP          string
	EnableVideo bool
	EnableAudio bool
	ScanTimeout time.Duration
	CtrlTimeOut time.Duration
}

func getDefaultOption(option *Option) *Option {
	if option == nil {
		return defaultOption
	}
	if option.CtrlTimeOut == 0 {
		option.CtrlTimeOut = defaultOption.CtrlTimeOut
	}
	if option.ScanTimeout == 0 {
		option.ScanTimeout = defaultOption.ScanTimeout
	}
	return option
}
