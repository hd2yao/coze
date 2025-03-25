package main

import (
	"time"

	"gorm.io/gorm"
)

// Alarm 闹钟数据库模型
type Alarm struct {
	gorm.Model
	Time        time.Time `gorm:"index;type:datetime;default:CURRENT_TIMESTAMP"` // 闹钟时间
	Schedule    string    `gorm:"type:varchar(100)"`                             // cron表达式
	Message     string    `gorm:"type:varchar(255)"`                             // 提醒消息
	AlarmType   string    `gorm:"type:varchar(20)"`                              // 闹钟类型：one_time/recurring
	IsActive    bool      `gorm:"default:true"`                                  // 是否激活
	CronEntryID int       `gorm:"type:int"`                                      // cron任务ID
}

// TableName 指定表名
func (alarm *Alarm) TableName() string {
	return "alarms_clock"
}
