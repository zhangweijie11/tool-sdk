package models

import "gorm.io/datatypes"

type Result struct {
	ID       uint           `gorm:"column:id;primarykey" json:"id"`
	UUID     string         `gorm:"column:uuid;size:37;uniqueIndex;comment:结果uuid" json:"uuid"`
	WorkUUID string         `gorm:"column:work_uuid;size:37;index;comment:总任务uuid" json:"work_uuid"`
	TaskUUID string         `gorm:"column:work_uuid;size:37;index;comment:子任务uuid" json:"task_uuid"`
	IP       string         `gorm:"column:ip;size:150;index;not null;comment:IP" json:"IP"`
	Port     uint16         `gorm:"column:port;index;not null;comment:开放端口" json:"port"`
	Status   string         `gorm:"column:status;size:64;index;not null;comment:端口状态" json:"status"`
	Protocol string         `gorm:"column:protocol;size:64;index;not null;default='';comment:协议" json:"protocol"`
	Service  string         `gorm:"column:service;size:128;index;not null;default='';comment:端口服务" json:"service"`
	Product  string         `gorm:"column:product;size:128;index;not null;default='';comment:组件产品" json:"product"`
	Version  string         `gorm:"column:version;size:64;index;not null;default='';comment:版本信息" json:"version"`
	Device   string         `gorm:"column:device;size:128;index;not null;default='';comment:设备类型" json:"device"`
	Extra    datatypes.JSON `gorm:"column:extra;type:json,not null;default='{}';comment:全量数据" json:"extra"`
	ComplexBaseModel
}

func (Result) TableName() string {
	return "result"
}
