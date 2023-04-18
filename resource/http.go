package resource

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/databonfire/bonfire/ac"
	"github.com/databonfire/bonfire/filter"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gorm.io/gorm/schema"
)

type HTTPHandlerMiddleware func(http.HandlerFunc) http.HandlerFunc

func RegisterHTTPServer(srv *http.Server, opt *Option) func() {
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
	repoName := opt.Resource
	if opt.AuthPackage {
		repoName = "auth." + repoName
	}
	registerRepo(repoName, repo)
	svc := NewService(opt, repo)

	pathPrefix := "/"
	if opt.Parent != "" {
		pathPrefix += opt.Parent + "/{pid}/"
	}
	pathPrefix += opt.Resource
	r := srv.Route("/")
	r.GET(pathPrefix, AssembleHandler(listHTTPHandler(svc), opt.HTTPHandlerMiddlewares))
	r.GET(pathPrefix+"/{id}", AssembleHandler(showHTTPHandler(svc), opt.HTTPHandlerMiddlewares))
	r.POST(pathPrefix, AssembleHandler(createHTTPHandler(svc), opt.HTTPHandlerMiddlewares))
	r.POST(pathPrefix+"/{id}", AssembleHandler(updateHTTPHandler(svc), opt.HTTPHandlerMiddlewares))
	r.PUT(pathPrefix+"/{id}", AssembleHandler(updateHTTPHandler(svc), opt.HTTPHandlerMiddlewares))
	r.DELETE(pathPrefix+"/{id}", AssembleHandler(deleteHTTPHandler(svc), opt.HTTPHandlerMiddlewares))

	return cleanup
}

func AssembleHandler(f http.HandlerFunc, mws []HTTPHandlerMiddleware) http.HandlerFunc {
	for i := len(mws) - 1; i >= 0; i-- {
		f = mws[i](f)
	}
	return f
}

func listHTTPHandler(svc *Service) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var lr ListRequest
		if err := ctx.BindQuery(&lr); err != nil {
			return err
		}
		if lr.FilterJsonlized != "" {
			if err := json.Unmarshal([]byte(lr.FilterJsonlized), &lr.Filter); err != nil {
				return err
			}
		}
		//if err := ctx.Bind(&lr); err != nil {
		//	return err
		//}

		reply, err := ctx.Middleware(func(stdctx context.Context, req interface{}) (interface{}, error) {
			// Parent access control
			if svc.Option.Parent != "" {
				var r RecordRequest
				if err := ctx.BindVars(&r); err != nil {
					return nil, err
				}
				if _, err := allowActionResource(stdctx, ac.ActionShow, svc.Option.Parent, r.PID); err != nil {
					return nil, err
				}
				if lr.Filter == nil {
					lr.Filter = filter.Filter{}
				}
				// GET /posts/1/comments?filters={star:{gt:4}}
				(map[string]interface{})(lr.Filter)[schema.NamingStrategy{}.ColumnName("", svc.Option.ParentField)] = r.PID // {{post_id: 1, star: {gt: 4}}
			}

			data, total, err := svc.repo.List(stdctx, &lr)
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
		})(ctx, &lr)
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
			// Parent access control
			if svc.Option.Parent != "" {
				if _, err := allowActionResource(ctx, ac.ActionShow, svc.Option.Parent, r.PID); err != nil {
					return nil, err
				}
			}

			record, err := svc.repo.Find(ctx, r.ID)
			if err != nil {
				return nil, err
			}
			// Orphan resource access control
			if svc.Option.Parent == "" && (ctx.Value("acer") != nil && !ctx.Value("acer").(ac.AccessController).Allow(ctx.Value("author"), ac.ActionShow, svc.Option.Resource, record)) {
				return nil, ac.ErrPermissionDenied
			}
			return record, nil
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
		if err := Validate(record); err != nil {
			return err
		}

		reply, err := ctx.Middleware(func(stdctx context.Context, req interface{}) (interface{}, error) {
			// Parent access control
			if svc.Option.Parent != "" {
				var r RecordRequest
				if err := ctx.BindVars(&r); err != nil {
					return nil, err
				}
				if _, err := allowActionResource(stdctx, ac.ActionEdit, svc.Option.Parent, r.PID); err != nil {
					return nil, err
				}
				reflect.ValueOf(record).Elem().FieldByName(svc.Option.ParentField).Set(reflect.ValueOf(r.PID))
			}
			if author := stdctx.Value("author"); author != nil {
				if accessor, ok := author.(ac.Accessor); ok {
					reflect.ValueOf(record).Elem().FieldByName("CreatedBy").Set(reflect.ValueOf(accessor.GetID()))
				}
			}
			return record, svc.repo.Save(stdctx, record)
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

		reply, err := ctx.Middleware(func(stdctx context.Context, req interface{}) (interface{}, error) {
			// Parent access control
			if svc.Option.Parent != "" {
				if _, err := allowActionResource(stdctx, ac.ActionEdit, svc.Option.Parent, r.PID); err != nil {
					return nil, err
				}
			}

			record, err := svc.repo.Find(stdctx, r.ID)
			if err != nil {
				return nil, err
			}
			// Orphan resource access control
			if svc.Option.Parent == "" && (stdctx.Value("acer") != nil && !stdctx.Value("acer").(ac.AccessController).Allow(stdctx.Value("author"), ac.ActionEdit, svc.Option.Resource, record)) {
				return nil, ac.ErrPermissionDenied
			}

			if err = ctx.Bind(record); err != nil {
				return nil, err
			}
			if err = validate.Struct(record); err != nil {
				return nil, err
			}
			reflect.ValueOf(record).Elem().FieldByName("ID").Set(reflect.ValueOf(r.ID))
			return record, svc.repo.Save(stdctx, record)
		})(ctx, r)
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

		reply, err := ctx.Middleware(func(stdctx context.Context, req interface{}) (interface{}, error) {
			// Parent access control
			if svc.Option.Parent != "" {
				if _, err := allowActionResource(stdctx, ac.ActionEdit, svc.Option.Parent, r.PID); err != nil {
					return nil, err
				}
			}

			record, err := svc.repo.Find(stdctx, r.ID)
			if err != nil {
				return nil, err
			}
			// Orphan resource access control
			if svc.Option.Parent == "" && (stdctx.Value("acer") != nil && !stdctx.Value("acer").(ac.AccessController).Allow(stdctx.Value("author"), ac.ActionDelete, svc.Option.Resource, record)) {
				return nil, ac.ErrPermissionDenied
			}

			return nil, svc.repo.Delete(stdctx, r.ID)
		})(ctx, r.ID)
		if err != nil {
			return err
		}
		return ctx.Result(200, reply)
	}
}

var (
	ErrExceedMaxNestDepth = errors.New("the depth of nest resource limit exceed")
	ErrInvalidParentID    = errors.New("unknown parent id in nest resource or type not support")
)
