package models

import (
	"database/sql/driver"
	"fmt"
	"gitlab.example.com/zhangweijie/tool-sdk/global"
	"gitlab.example.com/zhangweijie/tool-sdk/global/utils"
	"strings"
	"time"
)

type ComplexBaseModel struct {
	CreateTime DBTime `gorm:"column:create_time;not null;autoCreateTime" json:"create_time"`
	UpdateTime DBTime `gorm:"column:update_time;not null;autoUpdateTime" json:"update_time"`
}

// DBTime 自定义数据库时间类型
type DBTime struct {
	time.Time
}

func (t *DBTime) UnmarshalJSON(data []byte) error {
	var timestr = utils.BytesToStr(data)
	if timestr == "null" {
		return nil
	}
	t1, err := time.Parse(global.TimeFormatSecond, strings.Trim(timestr, "\""))
	*t = DBTime{t1}
	return err
}

func (t *DBTime) MarshalJSON() ([]byte, error) {
	return utils.StrToBytes(fmt.Sprintf("\"%s\"", t.Format(global.TimeFormatSecond))), nil
}

func (t *DBTime) Value() (driver.Value, error) {
	var zero time.Time
	if t.Time.UnixNano() == zero.UnixNano() {
		return nil, nil
	}
	return t.Format(global.TimeFormatSecond), nil
}

func (t *DBTime) Scan(v interface{}) error {
	switch value := v.(type) {
	case string:
		n, _ := time.Parse(global.TimeFormatSecond, value)
		*t = DBTime{n}
	case time.Time:
		*t = DBTime{value}
	default:
		return fmt.Errorf("can not convert %v to timestamp", v)
	}
	return nil
}
