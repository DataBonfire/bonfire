package resource

import (
	"reflect"
	"strings"

	"google.golang.org/grpc"
)

func RegisterGRPCServer(s *grpc.Server, opt *Option) func() {
	if nested := strings.Split(opt.Resource, "."); len(nested) > 1 {
		if len(nested) > 2 {
			panic(ErrExceedMaxNestDepth)
		}
		opt.Parent, opt.Resource = nested[0], nested[1]
		if v := strings.Split(opt.Parent, ":"); len(v) == 2 {
			opt.Parent, opt.ParentField = v[0], v[1]
		} else {
			opt.ParentField = toWord(opt.Parent[:len(opt.Parent)-1]) + "ID"
		}
		if t, ok := reflect.TypeOf(opt.Model).Elem().FieldByName(opt.ParentField); !ok || t.Type.Name() != "uint" {
			panic(ErrInvalidParentID)
		}
	}

	registeredDataMtx.Lock()
	var (
		data    *Data
		cleanup func()
	)
	if data = registeredData[opt.DataConfig]; data == nil {
		var err error
		data, cleanup, err = NewData(opt.DataConfig, opt.Logger)
		if err != nil {
			panic(err)
		}
		registeredData[opt.DataConfig] = data
	}
	registeredDataMtx.Unlock()
	repo := NewRepo(data, opt.Resource, opt.Model, opt.Logger)
	registerRepo(opt.Resource, repo)
	svc := NewService(opt, repo)

	s.RegisterService(makeServiceDesc(opt), svc)

	return cleanup
}

func makeServiceDesc(opt *Option) *grpc.ServiceDesc {
	return &grpc.ServiceDesc{
		ServiceName: "api." + opt.Resource,
		HandlerType: (*Service)(nil),
		Methods:     []grpc.MethodDesc{},
		Streams:     []grpc.StreamDesc{},
	}
}
