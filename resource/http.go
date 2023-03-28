package resource

import (
	"context"
	"reflect"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterHTTPServer(s *http.Server, opt *Option) func() {
	repo, cleanup, err := NewRepo(opt.DataConfig, opt.Model, opt.Logger)
	if err != nil {
		panic(err)
	}
	registerRepo(opt.Resource, repo)
	svc := NewService(opt, repo)

	r := s.Route("/")
	r.GET("/"+opt.Resource, listHTTPHandler(svc))
	r.GET("/"+opt.Resource+"/{id}", showHTTPHandler(svc))
	r.POST("/"+opt.Resource, createHTTPHandler(svc))
	r.POST("/"+opt.Resource+"/{id}", updateHTTPHandler(svc))
	r.DELETE("/"+opt.Resource+"/{id}", deleteHTTPHandler(svc))

	return cleanup
}

func listHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r ListRequest
		if err := ctx.BindQuery(&r); err != nil {
			return err
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
					PerPage: r.PerPage,
					Paged:   r.Paged,
				},
			}, nil
		})(ctx, r)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

type RecordRequest struct {
	ID uint `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
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
