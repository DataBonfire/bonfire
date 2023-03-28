package resource

import "github.com/go-kratos/kratos/v2/log"

type Option struct {
	Resource string
	Model    interface{}

	DataConfig *DataConfig
	Logger     log.Logger
}
