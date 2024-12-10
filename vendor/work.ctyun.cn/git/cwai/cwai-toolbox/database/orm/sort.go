package orm

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/golang/glog"
	"gorm.io/gorm/schema"
)

const (
	// ErrOrderByNotSupported 不支持排序
	ErrOrderByNotSupported = "order by is not supported"

	// OrderByAsc 正序
	OrderByAsc = "ASC"
	// OrderByDesc 倒序
	OrderByDesc = "DESC"
)

var factory = &SortGeneratorFactory{
	cache: map[string]*SortGenerator{},
	lock:  &sync.Mutex{},
}

type SortGenerator struct {
	Mapping map[string]string
}

// GormColumnName GORM框架的列名规则
func GormColumnName(field reflect.StructField) string {
	tags := field.Tag
	if sortTag, ok := tags.Lookup("sort"); ok {
		return sortTag
	}

	columnName := ""
	if gormTag, ok := tags.Lookup("gorm"); ok {
		if gormTag == "-" {
			return ""
		}

		parts := strings.Split(gormTag, ";")
		for _, part := range parts {
			if strings.HasPrefix(part, "column:") {
				columnName = strings.TrimPrefix(part, "column:")
			}
		}
	}

	if columnName == "" {
		naming := schema.NamingStrategy{}
		columnName = naming.ColumnName("", field.Name)
	}
	return columnName
}

// Generate Automatically
func NewSortGenerator(st interface{}) *SortGenerator {
	mapping := map[string]string{}

	stVals := make([]reflect.Value, 1, 5)
	stVals[0] = reflect.ValueOf(st)

	for len(stVals) > 0 {
		stVal := stVals[0]
		stVals = stVals[1:]

		if stVal.Type().Kind() == reflect.Ptr {
			stVal = stVal.Elem()
		}
		T := stVal.Type()
		glog.Infof("deducing mapping from type %s", T.Name())

		for i := 0; i < T.NumField(); i++ {
			field := T.Field(i)
			if field.Anonymous {
				if field.Type.Kind() == reflect.Struct {
					glog.Infof("adding anonymous field: %s", field.Name)
					stVals = append(stVals, stVal.Field(i))
				} else {
					glog.Infof("skipping anonymous field: %s", field.Name)
				}
				continue
			}

			tags := field.Tag
			dbTag := GormColumnName(field)
			jsonTag := tags.Get("json")
			if jsonTag == "-" || jsonTag == "" {
				jsonTag = dbTag
			}

			if dbTag != "" && jsonTag != "" {
				if dbTagExisting, exists := mapping[jsonTag]; exists {
					glog.Fatalf("new entry %s=>%s collide with %s=>%s, aborting", jsonTag, dbTag, jsonTag, dbTagExisting)
					return nil
				}
				mapping[jsonTag] = dbTag
			}
			if dbTag != "" {
				mapping[dbTag] = dbTag
			}
		}
	}

	glog.Infof("built sort mapping with %d entries", len(mapping))
	return &SortGenerator{Mapping: mapping}
}

func (g SortGenerator) Translate(jsonTagName string) (sqlTagName string) {
	return g.Mapping[jsonTagName]
}

func (g SortGenerator) GenerateFields(sortString string) (sortFields string, err error) {
	trimmedSortString := strings.TrimSpace(sortString)
	if len(trimmedSortString) == 0 {
		return "", nil
	}

	// field;field,desc;field,asc;...
	parts := strings.Split(sortString, ";")
	if len(parts) == 0 {
		return "", nil
	}

	var queries []string
	for _, part := range parts {
		partParts := strings.SplitN(part, ",", 2)
		if len(partParts) == 0 {
			continue
		}

		dbFieldName, isTranslated := g.Mapping[partParts[0]]
		if !isTranslated {
			return "", fmt.Errorf("column %s is not supported sort", partParts[0])
		}

		// SplitN assures there's no 3 or 4 parts, only 0, 1 and 2.
		if len(partParts) == 1 {
			queries = append(queries, dbFieldName)
		} else {
			order := strings.ToUpper(partParts[1])
			if order != OrderByAsc && order != OrderByDesc {
				return "", fmt.Errorf("order %s is not supported", partParts[1])
			}
			queries = append(queries, dbFieldName+" "+order)
		}
	}

	if len(queries) == 0 {
		return "", nil
	}

	return strings.Join(queries, ", "), nil
}

// Generate 生成ORDER BY子句
func (g SortGenerator) Generate(sortString string) (string, error) {
	fields, err := g.GenerateFields(sortString)
	if err != nil {
		return "", err
	}
	if fields != "" {
		fields = "ORDER BY " + fields
	}
	return fields, nil
}

// Validate 校验排序字段
func (g SortGenerator) Validate(sortString string) error {
	_, err := g.GenerateFields(sortString)
	return err
}

// SortGeneratorFactory 缓存已解析过的SortGenerator，避免重复解析
type SortGeneratorFactory struct {
	cache map[string]*SortGenerator
	lock  sync.Locker
}

func (factory *SortGeneratorFactory) GetGenerator(m interface{}) (*SortGenerator, error) {
	factory.lock.Lock()
	defer factory.lock.Unlock()

	key := reflect.TypeOf(m).String()
	g, ok := factory.cache[key]
	if ok {
		return g, nil
	}

	g = NewSortGenerator(m)
	if g == nil {
		return nil, fmt.Errorf("failed to generate sort generator for %s", key)
	}
	factory.cache[key] = g
	return g, nil
}

// ValidateSort 根据模型m校验排序字段sort
// m 为gorm的Model, 将json的Tag映射为gorm的column
// sort 的格式为column1,asc;column2,desc;column3,...
func ValidateSort(m interface{}, sort string) error {
	validator, err := factory.GetGenerator(m)
	if err != nil {
		return err
	}
	return validator.Validate(sort)
}

func GenerateSortFields(m interface{}, sort string) (string, error) {
	validator, err := factory.GetGenerator(m)
	if err != nil {
		return "", err
	}
	return validator.GenerateFields(sort)
}
