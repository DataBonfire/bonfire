package data

import (
	"fmt"

	"github.com/databonfire/bonfire/ac/rbac"
	"github.com/databonfire/bonfire/examples/singleton/internal/conf"
	"github.com/databonfire/bonfire/filter"
	"github.com/databonfire/bonfire/resource"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

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
		err = fmt.Errorf("unsupported database driver:", c.Database.Driver)
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
	{Model: gorm.Model{ID: 1}, Actions: resource.StringSlice{"browse"}, Resource: "posts"},
	{Model: gorm.Model{ID: 2}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts"},
	{Model: gorm.Model{ID: 3}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts", Record: filter.Filter{"created_by": "S"}},
	{Model: gorm.Model{ID: 4}, Actions: resource.StringSlice{"create", "edit", "delete"}, Resource: "posts", Record: filter.Filter{"created_by": "U"}},
}
var seeds = []interface{}{
	[]*rbac.Role{
		{Model: gorm.Model{ID: 1}, Name: "editor", Type: "user", Permissions: []*rbac.Permission{
			perms[0],
			perms[3],
		}},
		{Model: gorm.Model{ID: 2}, Name: "editor_manager", Type: "user", Permissions: []*rbac.Permission{
			perms[0],
			perms[2],
			perms[3],
		}},
		{Model: gorm.Model{ID: 3}, Name: "admin", Type: "admin", Permissions: []*rbac.Permission{
			perms[0],
			perms[1],
		}},
	},
}

func seedDB(db *gorm.DB) error {
	for _, seed := range seeds {
		db.AutoMigrate(seed)
		if err := db.Save(seed).Error; err != nil {
			return err
		}
	}
	return nil
}
