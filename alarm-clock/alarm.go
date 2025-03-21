package main

import (
    "fmt"
    "time"

    "github.com/robfig/cron/v3"
)

// Alarm 结构体
type Alarm struct {
    ID      cron.EntryID
    Time    time.Time
    Repeat  string // "none", "daily", "weekly"
    Message string
}

// AlarmManager 管理闹钟
type AlarmManager struct {
    cronScheduler *cron.Cron
    alarms        map[cron.EntryID]Alarm
}

// NewAlarmManager 创建 AlarmManager
func NewAlarmManager() *AlarmManager {
    return &AlarmManager{
        cronScheduler: cron.New(cron.WithSeconds()),
        alarms:        make(map[cron.EntryID]Alarm),
    }
}

// AddOneTimeAlarm 添加单次提醒
func (am *AlarmManager) AddOneTimeAlarm(t time.Time, msg string) {
    id, _ := am.cronScheduler.AddFunc(t.Format("05 04 15 02 01 *"), func() {
        fmt.Println("[闹钟提醒]：", msg)
    })
    am.alarms[id] = Alarm{ID: id, Time: t, Repeat: "none", Message: msg}
}

// AddRepeatingAlarm 添加重复提醒
func (am *AlarmManager) AddRepeatingAlarm(schedule, msg string) {
    id, _ := am.cronScheduler.AddFunc(schedule, func() {
        fmt.Println("[重复闹钟提醒]：", msg)
    })
    am.alarms[id] = Alarm{ID: id, Repeat: schedule, Message: msg}
}

// RemoveAlarm 移除闹钟
func (am *AlarmManager) RemoveAlarm(id cron.EntryID) {
    am.cronScheduler.Remove(id)
    delete(am.alarms, id)
}

// ListAlarms 列出所有闹钟
func (am *AlarmManager) ListAlarms() {
    for _, alarm := range am.alarms {
        fmt.Printf("ID: %d, 时间: %v, 类型: %s, 内容: %s\n", alarm.ID, alarm.Time, alarm.Repeat, alarm.Message)
    }
}

// Start 启动调度器
func (am *AlarmManager) Start() {
    am.cronScheduler.Start()
}
