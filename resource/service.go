package resource

import (
	"context"
	"reflect"
	"strings"

	"github.com/databonfire/bonfire/ac"
)

type Service struct {
	*Option
	repo         Repo
	resource     interface{}
	resourceType reflect.Type
}

func NewService(opt *Option, repo Repo) *Service {
	return &Service{
		Option:       opt,
		repo:         repo,
		resource:     opt.Model,
		resourceType: reflect.TypeOf(opt.Model),
	}
}

func allowActionResource(ctx context.Context, act string, res string, id uint) (interface{}, error) {
	record, err := GetRepo(res).Find(ctx, id)
	if err != nil {
		return nil, err
	}
	acer, ok := ctx.Value("acer").(ac.AccessController)
	if ok && !acer.Allow(ctx.Value("author"), act, res, record) {
		return nil, ac.ErrPermissionDenied
	}
	return record, nil
}

func toWord(s string) string {
	var worlds []string
	for _, v := range strings.Split(s, "_") {
		worlds = append(worlds, strings.ToUpper(v[0:1])+v[1:])
	}
	return strings.Join(worlds, "")
}
