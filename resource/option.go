package resource

import "github.com/go-kratos/kratos/v2/log"

type Option struct {
	Parent                 string
	ParentField            string
	AuthPackage            bool
	Resource               string
	Model                  interface{}
	HTTPHandlerMiddlewares []HTTPHandlerMiddleware

	DataConfig *DataConfig
	Logger     log.Logger
}
