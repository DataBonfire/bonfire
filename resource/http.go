package resource

import (
	"context"
	"errors"
	"reflect"
	"strings"

	"github.com/go-kratos/kratos/v2/transport/http"
	"gorm.io/gorm/schema"
)

func RegisterHTTPServer(s *http.Server, opt *Option) func() {
	// TODO add support multiple resource nest and string pid
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
	repo, cleanup, err := NewRepo(opt.DataConfig, opt.Model, opt.Logger)
	if err != nil {
		panic(err)
	}
	registerRepo(opt.Resource, repo)
	svc := NewService(opt, repo)

	pathPrefix := "/"
	if opt.Parent != "" {
		pathPrefix += opt.Parent + "/{pid}/"
	}
	pathPrefix += opt.Resource
	r := s.Route("/")
	r.GET(pathPrefix, listHTTPHandler(svc))
	r.GET(pathPrefix+"/{id}", showHTTPHandler(svc))
	r.POST(pathPrefix, createHTTPHandler(svc))
	r.POST(pathPrefix+"/{id}", updateHTTPHandler(svc))
	r.DELETE(pathPrefix+"/{id}", deleteHTTPHandler(svc))

	return cleanup
}

func listHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var lr ListRequest
		if err := ctx.BindQuery(&lr); err != nil {
			return err
		}
		if svc.Option.Parent != "" {
			var r RecordRequest
			if err := ctx.BindVars(&r); err != nil {
				return err
			}
			if lr.Filter == nil {
				lr.Filter = &Filter{}
			}
			(map[string]interface{})(*lr.Filter)[schema.NamingStrategy{}.ColumnName("", svc.Option.ParentField)] = r.PID
		}

		reply, err := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			data, total, err := svc.repo.List(ctx, req.(*ListRequest))
			if err != nil {
				return nil, err
			}
			return &ListResponse{
				Data: data,
				Pagination: &Pagination{
					Total:   total,
					PerPage: lr.PerPage,
					Paged:   lr.Paged,
				},
			}, nil
		})(ctx, lr)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

type RecordRequest struct {
	PID uint `protobuf:"bytes,1,opt,name=pid,proto3" json:"pid,omitempty"`
	ID  uint `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
}

func showHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r RecordRequest
		if err := ctx.BindVars(&r); err != nil {
			return err
		}

		reply, err := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return svc.repo.Find(ctx, req.(uint))
		})(ctx, &r.ID)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

func createHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		record := reflect.New(svc.resourceType.Elem()).Interface()
		if err := ctx.Bind(record); err != nil {
			return err
		}
		if svc.Option.Parent != "" {
			var r RecordRequest
			if err := ctx.BindVars(&r); err != nil {
				return err
			}
			reflect.ValueOf(record).Elem().FieldByName(svc.Option.ParentField).Set(reflect.ValueOf(r.PID))
		}

		reply, err := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return record, svc.repo.Save(ctx, req)
		})(ctx, record)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

func updateHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r RecordRequest
		if err := ctx.BindVars(&r); err != nil {
			return err
		}

		// TODO req should be record
		reply, err := ctx.Middleware(func(stdctx context.Context, req interface{}) (interface{}, error) {
			record, err := svc.repo.Find(ctx, req.(uint))
			if err != nil {
				return nil, err
			}
			if err = ctx.Bind(record); err != nil {
				return nil, err
			}
			if !ctx.Value("author").(AC).Allow("edit", svc.Option.Resource, record) {
				return nil, ErrPermissionDenied
			}
			return record, svc.repo.Save(stdctx, record)
		})(ctx, r.ID)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

func deleteHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r RecordRequest
		if err := ctx.BindVars(&r); err != nil {
			return err
		}

		reply, err := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, svc.repo.Delete(ctx, r.ID)
		})(ctx, r.ID)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

func toWord(s string) string {
	var worlds []string
	for _, v := range strings.Split(s, "_") {
		worlds = append(worlds, strings.ToUpper(v[0:1])+v[1:])
	}
	return strings.Join(worlds, "")
}

var (
	ErrExceedMaxNestDepth = errors.New("the depth of nest resource limit exceed")
	ErrInvalidParentID    = errors.New("unknown parent id in nest resource or type not support")
)
