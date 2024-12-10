package clickhouse

import (
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse"
)

// ClickHouse
type ClickHouse struct {
	uri  string
	conn *sqlx.DB
}

// NewClickHouse 新建ClickHouse
func NewClickHouse(uri string) *ClickHouse {
	return &ClickHouse{
		uri: uri,
	}
}

// Connect 连接ClickHouse并缓存Connection
func (ch *ClickHouse) Connect() error {
	conn, err := Connect(ch.uri)
	if err != nil {
		return err
	}
	ch.conn = conn
	return nil
}

// GetClickHouse 获取连接
func (ch *ClickHouse) GetClickHouse() *sqlx.DB {
	return ch.conn
}

// Connect 连接ClickHouse
func Connect(uri string) (*sqlx.DB, error) {
	glog.Infof("dialing clickhouse server at %s", uri)
	c, err := sqlx.Open("clickhouse", uri)
	if err != nil {
		glog.Fatalf("failed to connect to clickhouse database: %s", err)
		return nil, err
	}

	var serverIsOnline bool
	err = c.Get(&serverIsOnline, `SELECT 1 AS ServerIsOnline`)
	if err != nil {
		glog.Fatalf("failed to ping clickhouse: %s", err)
		return nil, err
	}

	glog.Info("clickhouse connected.")
	return c, nil
}

var Default *ClickHouse

func InitDefault(uri string) error {
	ch := NewClickHouse(uri)
	if err := ch.Connect(); err != nil {
		return err
	}
	Default = ch
	return nil
}

// GetClickHouse 获取连接
func GetClickHouse() *sqlx.DB {
	if Default == nil {
		return nil
	}
	return Default.GetClickHouse()
}

// CHProxy Bug Circumvention
func IsClickHouseNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "404 page not found")
}

func Select(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	t := time.Now()
	err := db.Select(dest, query, args...)
	glog.Infof("CLICKHOUSE[%s]: %s, %v", time.Since(t), query, args)
	return err
}

func Get(db *sqlx.DB, dest interface{}, query string, args ...interface{}) error {
	t := time.Now()
	err := db.Get(dest, query, args...)
	glog.Infof("CLICKHOUSE[%s]: %s, %v", time.Since(t), query, args)
	return err
}
