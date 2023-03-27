package resource

import "reflect"

type Service struct {
	*Option
	repo         Repo
	resource     interface{}
	resourceType reflect.Type
}

func NewService(opt *Option, res interface{}, repo Repo) *Service {
	return &Service{
		Option:       opt,
		repo:         repo,
		resource:     res,
		resourceType: reflect.TypeOf(res),
	}
}
