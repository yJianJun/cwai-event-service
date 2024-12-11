package orm

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Config struct {
	URI             string        `mapstructure:"uri" json:"uri" yaml:"uri"`
	MaxIdleConns    int           `mapstructure:"maxidle" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxopen" json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"maxlifetime" json:"connMaxLifetime" yaml:"connMaxLifetime"`
	LogEnabled      bool          `mapstructure:"logEnabled" json:"logEnabled" yaml:"logEnabled"`
}

var DebugConfig = &Config{
	URI:             "sqlite3://db.sqlite3",
	MaxIdleConns:    1,
	MaxOpenConns:    1,
	ConnMaxLifetime: time.Hour,
	LogEnabled:      true,
}

// New 新建ORM连接对象
func New(conf *Config) (*gorm.DB, error) {
	uri := conf.URI
	var dialector gorm.Dialector
	if strings.HasPrefix(uri, "mysql://") {
		dialector = mysql.Open(uri[8:])
	} else if strings.HasPrefix(uri, "sqlite3://") {
		dialector = sqlite.Open(uri[10:])
	} else {
		dialector = mysql.Open(uri)
	}

	ormDB, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		glog.Errorf("failed to open gorm: config=%+v, error=%v", conf, err)
		return nil, err
	}
	db, err := ormDB.DB()
	if err != nil {
		glog.Errorf("failed to get database connection: uri=%v, error=%v", uri, err)
		return nil, err
	}

	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	if conf.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(conf.ConnMaxLifetime)
	}
	if conf.LogEnabled {
		ormDB = ormDB.Debug()
	}
	return ormDB, nil
}
