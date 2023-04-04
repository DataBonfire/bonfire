package data

import (
	"fmt"

	"github.com/databonfire/bonfire/ac/rbac"
	"github.com/databonfire/bonfire/examples/singleton/internal/conf"
	"github.com/databonfire/bonfire/filter"
	"github.com/databonfire/bonfire/resource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	var (
		db  *gorm.DB
		err error
	)
	switch c.Database.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	default:
		err = fmt.Errorf("unsupported database driver:%s", c.Database.Driver)
	}
	if err != nil {
		return nil, nil, err
	}

	if err = seedDB(db); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db}, cleanup, nil
}

var perms = []*rbac.Permission{
	{Model: resource.Model{ID: 1}, Actions: resource.StringSlice{"browse"}, Resource: "posts"},
	{Model: resource.Model{ID: 2}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts"},
	{Model: resource.Model{ID: 3}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts", Record: filter.Filter{"created_by": "S"}},
	{Model: resource.Model{ID: 4}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts", Record: filter.Filter{"created_by": "U"}},
}
var seeds = []interface{}{
	rbac.MakeRoles(&rbac.RoleTemplate{
		Type:       "admin",
		Name:       "admin",
		ActionsAll: []string{"posts"},
	}, &rbac.RoleTemplate{
		Type:                  "user",
		Name:                  "editor_manager",
		ActionsMyCreated:      []string{"posts.created_by"},
		ActionsMySubordinates: []string{"posts.created_by"},
	}, &rbac.RoleTemplate{
		Type:             "user",
		Name:             "editor",
		ActionsMyCreated: []string{"posts.created_by"},
	}),
}

func seedDB(db *gorm.DB) error {
	db.AutoMigrate(&rbac.Role{})
	for _, seed := range seeds {
		db.AutoMigrate(seed)
		if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(seed).Error; err != nil {
			return err
		}
	}
	return nil
}
