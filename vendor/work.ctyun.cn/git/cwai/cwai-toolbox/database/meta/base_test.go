package meta

import (
	"testing"

	"work.ctyun.cn/git/cwai/cwai-toolbox/database/mysql"
	"work.ctyun.cn/git/cwai/cwai-toolbox/database/orm"

	"github.com/stretchr/testify/assert"
)

type User struct {
	Meta
	Name string
}

func TestMeta(t *testing.T) {
	db, err := orm.New(orm.DebugConfig)
	assert.Nil(t, err)

	user, u := User{Name: "toolbox"}, User{}
	user.SetDefault()
	db.Migrator().DropTable(&User{})
	db.AutoMigrate(&User{})
	db.Save(&user)

	err = db.Model(&u).Where("name = ?", user.Name).Last(&u).Error
	assert.Nil(t, err)
	assert.Equal(t, user, u)
	t.Logf("got: %+v", u)

	err = db.Delete(&user).Error
	assert.Nil(t, err)
	assert.True(t, user.IsDeleted())

	err = db.First(&User{}).Error
	assert.True(t, mysql.NotFound(err))

	u = User{}
	err = db.Unscoped().First(&u).Error
	assert.Nil(t, err)
	assert.True(t, user.IsDeleted())
	t.Logf("deleted: %+v", u)
}
