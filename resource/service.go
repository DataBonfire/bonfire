package resource

import "reflect"

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
