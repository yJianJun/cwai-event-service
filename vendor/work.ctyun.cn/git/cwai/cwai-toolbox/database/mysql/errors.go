package mysql

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrNoRowAffected = errors.New("no row affected")
)

// NotFound: is not found error?
func NotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// NotAffected: is no row affected error?
func NotAffected(err error) bool {
	return err == ErrNoRowAffected
}

// DuplicateEntry: is duplicated entry error?
func DuplicateEntry(err error) bool {
	// Ref: https://github.com/go-gorm/gorm/issues/1718
	if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		return true
	}
	return false
}

// DetectErrors: detects error from db after sql execution
func DetectErrors(db *gorm.DB) error {
	if db.Error == nil && db.RowsAffected == 0 {
		return ErrNoRowAffected
	}
	return db.Error
}
