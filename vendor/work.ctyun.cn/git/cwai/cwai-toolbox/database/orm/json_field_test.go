package orm

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/suite"
)

type SomeField struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func (field SomeField) Value() (driver.Value, error) {
	return Value(field)
}

func (field *SomeField) Scan(value interface{}) error {
	return Scan(field, value)
}

type SomeModel struct {
	gorm.Model
	Some SomeField `gorm:"type:text"`
}

type JSONFieldTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

func (suite *JSONFieldTestSuite) SetupTest() {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	db.LogMode(true)
	suite.DB = db
	suite.Nil(err)
	suite.DB.CreateTable(&SomeModel{})
}

func (suite *JSONFieldTestSuite) TearDownTest() {
	suite.DB.DropTableIfExists(&SomeModel{})
}

func TestJSONField(t *testing.T) {
	suite.Run(t, new(JSONFieldTestSuite))
}

func (suite *JSONFieldTestSuite) TestJSON() {
	some := SomeModel{
		Some: SomeField{
			A: 1,
			B: "abc",
		},
	}
	err := suite.DB.Create(&some).Error
	suite.Nil(err)

	some2 := SomeModel{}
	err = suite.DB.First(&some2).Error
	suite.Nil(err)
	suite.Equal(1, some2.Some.A)
	suite.Equal("abc", some2.Some.B)

	v, err := some.Some.Value()
	fmt.Printf("%s\n", v.(string))

}
