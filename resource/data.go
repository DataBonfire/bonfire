package resource

import (
	"context"
	"fmt"
	"gorm.io/gorm/clause"
	"reflect"
	"sync"

	"github.com/databonfire/bonfire/ac"
	"github.com/databonfire/bonfire/filter"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	registeredData    = map[*DataConfig]*Data{}
	registeredDataMtx sync.Mutex
)

type Data struct {
	db *gorm.DB
}

func NewData(c *DataConfig, logger log.Logger) (*Data, func(), error) {
	var (
		db  *gorm.DB
		err error
	)
	switch c.Database.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{FullSaveAssociations: true})
	default:
		err = fmt.Errorf("unsupported database driver:%s", c.Database.Driver)
	}
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db.Debug()}, cleanup, nil
}

type Repo interface {
	DB() *gorm.DB
	List(context.Context, *ListRequest) ([]interface{}, int64, error)
	Find(context.Context, uint) (interface{}, error)
	Save(context.Context, interface{}) error
	Delete(context.Context, uint) error
}

type repo struct {
	data      *Data
	resource  string
	model     interface{}
	modelType reflect.Type
	log       *log.Helper
}

func NewRepo(data *Data, resource string, model interface{}, logger log.Logger) Repo {
	if err := data.db.AutoMigrate(model); err != nil {
		panic(err)
	}
	return &repo{
		data:      data,
		resource:  resource,
		model:     model,
		modelType: reflect.TypeOf(model),
		log:       log.NewHelper(logger),
	}
}

func (r *repo) DB() *gorm.DB {
	return r.data.db
}

func (r *repo) List(ctx context.Context, lr *ListRequest) ([]interface{}, int64, error) {
	var (
		total int64
		data  = reflect.New(reflect.MakeSlice(reflect.SliceOf(r.modelType), 0, 0).Type())
		//errs  = make(chan error, 1)
	)
	rootDB := r.data.db.WithContext(ctx).Preload(clause.Associations)
	db, err := filter.GormFilter(rootDB, lr.Filter)
	if err != nil {
		return nil, 0, err
	}
	if acer := ctx.Value("acer"); acer != nil {
		if filters := acer.(ac.AccessController).Filters(ctx.Value("author"), ac.ActionBrowse, r.resource); filters != nil {
			acDB, err := filter.GormFilter(rootDB, filters...)
			if err != nil {
				return nil, 0, err
			}
			db.Where(acDB)
		}
	}
	if err = db.Model(r.model).Count(&total).Offset(int(lr.Paged * lr.PerPage)).Limit(int(lr.PerPage)).Find(data.Interface()).Error; err != nil {
		return nil, 0, err
	}

	//db := r.data.db
	//go func() {
	//	for _, v := range lr.Sorts {
	//		db.Order(fmt.Sprintf("%s %s", v.By, v.Order))
	//	}
	//	if lr.Paged > 0 {
	//		db.Offset(int(lr.Paged * lr.PerPage))
	//	}
	//	if lr.PerPage > 0 {
	//		db.Limit(int(lr.PerPage))
	//	}
	//	tx := db.Find(data.Interface())
	//	errs <- tx.Error
	//}()
	//go func() {
	//	tx := db.Model(r.model).Count(&total)
	//	errs <- tx.Error
	//}()

	//for i := 0; i < 2; i++ {
	//	if err := <-errs; err != nil {
	//		return nil, 0, err
	//	}
	//}

	var list []interface{}
	for i := 0; i < data.Elem().Len(); i++ {
		list = append(list, data.Elem().Index(i).Interface())
	}
	return list, total, nil
}

func (r *repo) Find(ctx context.Context, id uint) (interface{}, error) {
	dest := reflect.New(r.modelType)
	tx := r.data.db.Preload(clause.Associations).First(dest.Interface(), id)
	return dest.Elem().Interface(), tx.Error
}

func (r *repo) Save(ctx context.Context, record interface{}) error {
	return r.data.db.Save(record).Error
}

func (r *repo) Delete(ctx context.Context, id uint) error {
	return r.data.db.Delete(r.model, id).Error
}
