package mysql

import (
	"context"

	"work.ctyun.cn/git/cwai/cwai-toolbox/database/orm"

	"github.com/golang/glog"
	"gorm.io/gorm"
)

var database *gorm.DB

// Init: initialize mysql
func Init(c *orm.Config) (err error) {
	glog.Infof("initializing mysql: %+v", c)
	database, err = orm.New(c)
	if err != nil {
		return err
	}
	return nil
}

// GetDB: get database
func GetDB() *gorm.DB {
	return database
}

// GetDB: get database with context
func GetDBWithContext(ctx context.Context) *gorm.DB {
	return database.WithContext(ctx)
}
