package resource

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Repo interface {
	List(context.Context, *ListRequest) ([]interface{}, int64, error)
	Find(context.Context, uint) (interface{}, error)
	Save(context.Context, interface{}) error
	Delete(context.Context, uint) error
}

type repo struct {
	db        *gorm.DB
	model     interface{}
	modelType reflect.Type
}

func NewRepo(c *Config, model interface{}, logger log.Logger) (Repo, func(), error) {
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
	if err = db.AutoMigrate(model); err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &repo{
		db:        db.Debug(),
		model:     model,
		modelType: reflect.TypeOf(model),
	}, cleanup, nil
}

func (r *repo) List(ctx context.Context, lr *ListRequest) ([]interface{}, int64, error) {
	var (
		total int64
		data  = reflect.New(reflect.MakeSlice(reflect.SliceOf(r.modelType), 0, 0).Type())
		errs  = make(chan error, 1)
	)
	go func() {
		tx := r.db.Find(data.Interface())
		errs <- tx.Error
	}()
	go func() {
		tx := r.db.Model(r.model).Count(&total)
		errs <- tx.Error
	}()

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			return nil, 0, err
		}
	}

	var list []interface{}
	for i := 0; i < data.Elem().Len(); i++ {
		list = append(list, data.Elem().Index(i).Interface())
	}
	return list, total, nil
}

func (r *repo) Find(ctx context.Context, id uint) (interface{}, error) {
	dest := reflect.New(r.modelType).Interface()
	tx := r.db.First(dest, id)
	return dest, tx.Error
}

func (r *repo) Save(ctx context.Context, record interface{}) error {
	return r.db.Save(record).Error
}

func (r *repo) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(r.model, id).Error
}
