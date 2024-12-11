package meta

import (
	"encoding/json"
	"time"

	"gorm.io/plugin/soft_delete"
)

/**
 * DateTime
 */
type DateTime int64

func NewDateTime(t time.Time) DateTime {
	return DateTime(t.Unix())
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	t := d.AsTime()
	formatted := t.Local().Format(time.RFC3339) // 2020-06-18T16:19:27+08:00
	return json.Marshal(formatted)
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	layout := `"` + time.RFC3339 + `"`
	if parsed, err := time.Parse(layout, string(data)); err != nil {
		return err
	} else {
		*t = DateTime(parsed.Unix())
		return nil
	}
}

func (d DateTime) String() string {
	return d.AsTime().Local().Format(time.RFC3339)
}

func (d DateTime) AsTime() time.Time {
	return time.Unix(int64(d), 0)
}

func Now() DateTime {
	return NewDateTime(time.Now())
}

/**
 * Meta
 */
// Meta: 定义了数据库表的基本信息，例如ID与创建时间等
type Meta struct {
	ID         int32                 `gorm:"primaryKey" json:"id"`
	Creator    string                `json:"creator"`
	CreateTime DateTime              `json:"createTime,omitempty"`
	UpdateTime DateTime              `json:"updateTime,omitempty"`
	DeleteTime soft_delete.DeletedAt `gorm:"softDelete,default:null" json:"deleteTime,omitempty"`
}

func (m *Meta) SetDefault() {
	now := Now()
	m.CreateTime = now
	m.UpdateTime = now
}

func (m *Meta) SetUpdateTime() {
	now := Now()
	m.UpdateTime = now
}

func (m *Meta) IsDeleted() bool {
	return m.DeleteTime > 0
}
