package main

type AlarmDao struct{}

func NewAlarmDao() *AlarmDao {
	return &AlarmDao{}
}

// CreateAlarm 创建新的闹钟记录
func (ad *AlarmDao) CreateAlarm(alarmModel *Alarm) error {
	return DBMaster().Create(alarmModel).Error
}

// GetAlarmByID 根据ID获取闹钟
func (ad *AlarmDao) GetAlarmByID(id uint) (*Alarm, error) {
	var alarmModel Alarm
	err := DB().First(&alarmModel, id).Error
	if err != nil {
		return nil, err
	}
	return &alarmModel, nil
}

// GetActiveAlarms 获取所有活动的闹钟
func (ad *AlarmDao) GetActiveAlarms() ([]Alarm, error) {
	var alarms []Alarm
	err := DB().Where("is_active = ?", true).Find(&alarms).Error
	if err != nil {
		return nil, err
	}
	return alarms, nil
}

// UpdateAlarm 更新闹钟信息
func (ad *AlarmDao) UpdateAlarm(alarmModel *Alarm) error {
	return DBMaster().Save(alarmModel).Error
}

// DeactivateAlarm 停用闹钟
func (ad *AlarmDao) DeactivateAlarm(id uint) error {
	return DBMaster().Model(&Alarm{}).Where("id = ?", id).Update("is_active", false).Error
}

// DeleteAlarm 删除闹钟
func (ad *AlarmDao) DeleteAlarm(id uint) error {
	return DBMaster().Delete(&Alarm{}, id).Error
}

// UpdateCronEntryID 更新闹钟的Cron任务ID
func (ad *AlarmDao) UpdateCronEntryID(alarmID uint, cronEntryID int) error {
	return DBMaster().Model(&Alarm{}).Where("id = ?", alarmID).Update("cron_entry_id", cronEntryID).Error
}
