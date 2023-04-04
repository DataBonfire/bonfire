package data

import (
	"fmt"

	"github.com/databonfire/bonfire/ac/rbac"
	"github.com/databonfire/bonfire/examples/singleton/internal/conf"
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

var (
	roles = []*rbac.RoleTemplate{
		{
			Type:       "admin",
			Name:       "admin",
			ActionsAll: []string{"posts"},
		},
		{
			Type:                "user",
			Name:                "editor_manager",
			ActionsUID:          []string{"posts.created_by"},
			ActionsSubordinates: []string{"posts.created_by"},
		},
		{
			Type:       "user",
			Name:       "editor",
			ActionsUID: []string{"posts.created_by"},
		},
	}
	seeds = []interface{}{
		rbac.MakeRoles(roles),
	}
)

func seedDB(db *gorm.DB) error {
	db.AutoMigrate(&rbac.Role{})
	for _, seed := range seeds {
		db.AutoMigrate(seed)
		if err := db.Save(seed).Error; err != nil {
			return err
		}
	}
	return nil
}
