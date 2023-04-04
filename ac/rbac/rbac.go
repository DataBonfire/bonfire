package rbac

import (
	"sync"
	"time"

	"github.com/databonfire/bonfire/ac"
	"github.com/databonfire/bonfire/filter"
	"github.com/databonfire/bonfire/resource"
	"github.com/go-kratos/kratos/v2/log"
)

type RBAC struct {
	repo          resource.Repo
	roleRegisters []*Role
	mtx           sync.RWMutex
	logger        *log.Helper
}

var rbac *RBAC

func newAC(repo resource.Repo, logger log.Logger) ac.AccessController {
	if rbac == nil {
		rbac = &RBAC{
			repo:   repo,
			logger: log.NewHelper(logger),
		}
		go func() {
			if err := rbac.refreshRegister(); err != nil {
				panic(err)
			}
			for range time.Tick(time.Minute * 5) {
				if err := rbac.refreshRegister(); err != nil {
					rbac.logger.Error(err)
				}
			}
		}()
	}
	return rbac
}

func (ac *RBAC) refreshRegister() error {
	if ac.roleRegisters == nil {
		ac.mtx.Lock()
		defer ac.mtx.Unlock()
	}
	if ac.repo == nil {
		for range time.Tick(time.Second) {
			if v := resource.GetRepo("roles"); v != nil {
				ac.repo = v.(resource.Repo)
				break
			}
		}
	}
	var roles []*Role
	if err := ac.repo.DB().Preload("Permissions").Find(&roles).Error; err != nil {
		return err
	}
	ac.roleRegisters = roles
	return nil
}

func contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (ac *RBAC) Allow(a interface{}, act string, res string, re interface{}) bool {
	accessor := accessorOrVisitor(a)
	ac.mtx.RLock()
	defer ac.mtx.RUnlock()
	for _, role := range ac.roleRegisters {
		if !contains(accessor.GetRoles(), role.Name) {
			continue
		}
		for _, perm := range role.Permissions {
			if !contains(perm.Actions, act) || perm.Resource != res {
				continue
			}

			// Not record lvl
			// No record perm
			if re == nil || perm.Record == nil {
				return true
			}

			pa := applyAccessor(accessor, perm.Record)
			if pa != nil && pa.Match(re) {
				return true
			}
		}
	}
	return false
}

func (ac *RBAC) Filters(a interface{}, act, res string) []filter.Filter {
	accessor := accessorOrVisitor(a)
	ac.mtx.RLock()
	defer ac.mtx.RUnlock()
	var filters []filter.Filter
	for _, role := range ac.roleRegisters {
		if !contains(accessor.GetRoles(), role.Name) {
			continue
		}
		for _, perm := range role.Permissions {
			if !contains(perm.Actions, act) || perm.Resource != res {
				continue
			}

			// No record lvl perm which means public to the role
			if perm.Record == nil {
				return nil
			}
			v := applyAccessor(accessor, perm.Record)
			if v == nil {
				continue
			}
			filters = append(filters, v)
		}
	}
	if len(filters) == 0 {
		filters = append(filters, map[string]interface{}{"created_by": -1})
	}

	return filters
}

func (ac *RBAC) Permissions(a interface{}) []*Permission {
	accessor := accessorOrVisitor(a)
	ac.mtx.RLock()
	defer ac.mtx.RUnlock()
	var perms []*Permission
	for _, role := range ac.roleRegisters {
		if contains(accessor.GetRoles(), role.Name) {
			for _, v := range role.Permissions {
				perm := &Permission{
					Actions:  v.Actions,
					Resource: v.Resource,
					Record:   applyAccessor(accessor, v.Record),
				}
				perms = append(perms, perm)
			}
		}
	}
	return perms
}

func accessorOrVisitor(v interface{}) Accessor {
	if v == nil {
		return &visitor
	}
	return v.(Accessor)
}

func applyAccessor(a Accessor, f filter.Filter) filter.Filter {
	filter := map[string]interface{}{}
	for k, v := range f {
		switch v {
		case "U":
			filter[k] = a.GetID()
		case "G":
			groups := a.GetGroups()
			switch len(groups) {
			case 0:
				return nil
				filter[k] = -1 // No group which means filter invalid
			case 1:
				filter[k] = groups[0]
			default:
				filter[k] = groups
			}
		case "S":
			subs := a.GetSubordinates()
			switch len(subs) {
			case 0:
				return nil
				filter[k] = -1 // No sub which means filter invalid
			case 1:
				filter[k] = subs[0]
			default:
				filter[k] = subs
			}
		default:
			filter[k] = v
		}
	}
	return filter
}
