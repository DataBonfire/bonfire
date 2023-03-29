package resource

import "github.com/go-kratos/kratos/v2/log"

type Option struct {
	Parent      string
	ParentField string
	Resource    string
	Model       interface{}

	DataConfig *DataConfig
	Logger     log.Logger
}
