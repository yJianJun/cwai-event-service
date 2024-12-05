package utils

import (
	"github.com/go-sql-driver/mysql"
	"strings"
)

func GetOffsetLimit(pageSize, pageNum int) (int, int) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNum <= 0 {
		pageNum = 1
	}
	return (pageNum - 1) * pageSize, pageSize
}

func GetPageCount(count, pageSize int) int {
	if pageSize == 0 {
		pageSize = 10
	}
	pages := count / pageSize
	if count%pageSize != 0 {
		pages++
	}
	return pages
}

func EscapeLike(str string) string {
	if str == "" {
		return str
	}

	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "_", "\\_", -1)
	str = strings.Replace(str, "%", "\\%", -1)

	return "%" + str + "%"
}

func IsDuplicateKeyError(err error) bool {
	if errMySQL, ok := err.(*mysql.MySQLError); ok && errMySQL.Number == 1062 {
		return true
	}
	return false
}

func IsRecordNotFoundError(err error) bool {
	return err.Error() == "record not found"
}

func IsMySQLError(err error) (bool, *mysql.MySQLError) {
	if errMySQL, ok := err.(*mysql.MySQLError); ok {
		return true, errMySQL
	}
	return false, nil
}