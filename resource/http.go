package resource

import (
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterHTTPServer(conf *Config, logger log.Logger, s *http.Server, opt *Option) func() {
	repo, cleanup, err := NewRepo(conf, opt.Model, logger)
	if err != nil {
		panic(err)
	}
	registerRepo(opt.Resource, repo)
	svc := NewService(opt, opt.Model, repo)

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
		data, total, err := svc.repo.List(ctx, &r)
		if err != nil {
			return err
		}
		reply := &ListResponse{
			Data: data,
			Pagination: &Pagination{
				Total:   total,
				PerPage: r.PerPage,
				Paged:   r.Paged,
			},
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
		reply, err := svc.repo.Find(ctx, r.ID)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

func createHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		record := reflect.New(svc.resourceType)
		if err := ctx.Bind(record); err != nil {
			return err
		}
		if err := svc.repo.Save(ctx, record); err != nil {
			return err
		}
		return ctx.Result(200, record)
	}
}

func updateHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r RecordRequest
		if err := ctx.BindVars(&r); err != nil {
			return err
		}
		record, err := svc.repo.Find(ctx, r.ID)
		if err != nil {
			return err
		}
		ctx.Value("user")(*biz.User).HasPermission("edit", record)
		if err = ctx.Bind(record); err != nil {
			return err
		}
		if err = svc.repo.Save(ctx, record); err != nil {
			return err
		}
		return ctx.Result(200, record)
	}
}

func deleteHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var r RecordRequest
		if err := ctx.BindVars(&r); err != nil {
			return err
		}
		if err := svc.repo.Delete(ctx, r.ID); err != nil {
			return err
		}
		return ctx.Result(200, nil)
	}
}
