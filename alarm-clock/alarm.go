package main

import (
    "fmt"
    "time"

    "github.com/robfig/cron/v3"
)

// AlarmManager 管理闹钟
type AlarmManager struct {
    cronScheduler *cron.Cron
    alarmDao      *AlarmDao
}

// NewAlarmManager 创建 AlarmManager
func NewAlarmManager() *AlarmManager {
    return &AlarmManager{
        cronScheduler: cron.New(cron.WithSeconds()),
        alarmDao:      NewAlarmDao(),
    }
}

// Start 启动调度器并恢复已存在的闹钟
func (am *AlarmManager) Start() error {
    // 启动cron调度器
    am.cronScheduler.Start()

    // 从数据库加载活动的闹钟
    alarms, err := am.alarmDao.GetActiveAlarms()
    if err != nil {
        return fmt.Errorf("加载闹钟失败: %v", err)
    }

    // 重新调度所有活动的闹钟
    for _, alarm := range alarms {
        if alarm.AlarmType == "one_time" {
            if alarm.Time.After(time.Now()) {
                am.scheduleOneTimeAlarm(&alarm)
            } else {
                // 如果一次性闹钟已过期，将其设置为非活动
                am.alarmDao.DeactivateAlarm(alarm.ID)
            }
        } else if alarm.AlarmType == "recurring" {
            am.scheduleRecurringAlarm(&alarm)
        }
    }

    return nil
}

// Stop 停止调度器
func (am *AlarmManager) Stop() {
    am.cronScheduler.Stop()
}

// AddOneTimeAlarm 添加单次提醒
func (am *AlarmManager) AddOneTimeAlarm(t time.Time, msg string) error {
    // 创建闹钟记录
    alarm := &Alarm{
        Time:      t,
        Message:   msg,
        AlarmType: "one_time",
        IsActive:  true,
    }

    // 保存到数据库
    if err := am.alarmDao.CreateAlarm(alarm); err != nil {
        return fmt.Errorf("保存闹钟失败: %v", err)
    }

    // 调度闹钟
    return am.scheduleOneTimeAlarm(alarm)
}

// AddRepeatingAlarm 添加重复提醒
func (am *AlarmManager) AddRepeatingAlarm(schedule, msg string) error {
    // 创建闹钟记录
    alarm := &Alarm{
        Schedule:  schedule,
        Message:   msg,
        AlarmType: "recurring",
        IsActive:  true,
    }

    // 保存到数据库
    if err := am.alarmDao.CreateAlarm(alarm); err != nil {
        return fmt.Errorf("保存闹钟失败: %v", err)
    }

    // 调度闹钟
    return am.scheduleRecurringAlarm(alarm)
}

// RemoveAlarm 移除闹钟
func (am *AlarmManager) RemoveAlarm(id uint) error {
    // 获取闹钟信息
    alarm, err := am.alarmDao.GetAlarmByID(id)
    if err != nil {
        return fmt.Errorf("获取闹钟失败: %v", err)
    }

    // 停止cron任务
    if alarm.CronEntryID != 0 {
        am.cronScheduler.Remove(cron.EntryID(alarm.CronEntryID))
    }

    // 从数据库中删除
    return am.alarmDao.DeleteAlarm(id)
}

// ListAlarms 列出所有闹钟
func (am *AlarmManager) ListAlarms() ([]Alarm, error) {
    return am.alarmDao.GetActiveAlarms()
}

// 内部方法：调度一次性闹钟
func (am *AlarmManager) scheduleOneTimeAlarm(alarm *Alarm) error {
    entryID, err := am.cronScheduler.AddFunc(alarm.Time.Format("05 04 15 02 01 *"), func() {
        fmt.Printf("[闹钟提醒] ID:%d - %s\n", alarm.ID, alarm.Message)
        // 提醒后将闹钟设置为非活动
        am.alarmDao.DeactivateAlarm(alarm.ID)
    })

    if err != nil {
        return fmt.Errorf("调度闹钟失败: %v", err)
    }

    // 更新CronEntryID
    return am.alarmDao.UpdateCronEntryID(alarm.ID, int(entryID))
}

// 内部方法：调度重复性闹钟
func (am *AlarmManager) scheduleRecurringAlarm(alarm *Alarm) error {
    entryID, err := am.cronScheduler.AddFunc(alarm.Schedule, func() {
        fmt.Printf("[重复闹钟提醒] ID:%d - %s\n", alarm.ID, alarm.Message)
    })

    if err != nil {
        return fmt.Errorf("调度闹钟失败: %v", err)
    }

    // 更新CronEntryID
    return am.alarmDao.UpdateCronEntryID(alarm.ID, int(entryID))
}
