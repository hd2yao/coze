package main

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/robfig/cron/v3"
)

// AlarmRequest 表示创建闹钟的请求
type AlarmRequest struct {
    Time     string `json:"time,omitempty"`     // 格式: "2006-01-02 15:04:05"
    Schedule string `json:"schedule,omitempty"` // cron格式的调度表达式
    Message  string `json:"message"`            // 闹钟消息
}

// AlarmResponse 表示API响应
type AlarmResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// AlarmAPI API 结构体处理HTTP请求
type AlarmAPI struct {
    manager *AlarmManager
    router  *gin.Engine
}

// NewAlarmAPI 创建新的API处理器
func NewAlarmAPI(manager *AlarmManager) *AlarmAPI {
    // 创建默认的Gin路由器
    router := gin.Default()

    api := &AlarmAPI{
        manager: manager,
        router:  router,
    }

    // 设置路由
    api.setupRoutes()

    return api
}

// 设置路由
func (api *AlarmAPI) setupRoutes() {
    // 闹钟相关接口
    alarmGroup := api.router.Group("/alarm")
    {
        alarmGroup.POST("/add", api.handleAddAlarm)
        alarmGroup.DELETE("/remove", api.handleRemoveAlarm)
        alarmGroup.GET("/list", api.handleListAlarms)
    }
}

// StartServer 启动HTTP服务器
func (api *AlarmAPI) StartServer(port string) {
    // 启动服务器
    api.router.Run(":" + port)
}

// 处理添加闹钟的请求
func (api *AlarmAPI) handleAddAlarm(c *gin.Context) {
    var req AlarmRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, AlarmResponse{
            Success: false,
            Message: "请求解析失败: " + err.Error(),
        })
        return
    }

    if req.Message == "" {
        c.JSON(http.StatusBadRequest, AlarmResponse{
            Success: false,
            Message: "闹钟消息不能为空",
        })
        return
    }

    // 根据请求类型添加不同类型的闹钟
    if req.Time != "" {
        // 添加一次性闹钟
        t, err := time.Parse("2006-01-02 15:04:05", req.Time)
        if err != nil {
            c.JSON(http.StatusBadRequest, AlarmResponse{
                Success: false,
                Message: "时间格式无效: " + err.Error(),
            })
            return
        }

        if t.Before(time.Now()) {
            c.JSON(http.StatusBadRequest, AlarmResponse{
                Success: false,
                Message: "闹钟时间不能设置为过去时间",
            })
            return
        }

        api.manager.AddOneTimeAlarm(t, req.Message)
        c.JSON(http.StatusOK, AlarmResponse{
            Success: true,
            Message: "一次性闹钟添加成功",
            Data: map[string]interface{}{
                "time":    t,
                "message": req.Message,
            },
        })
    } else if req.Schedule != "" {
        // 添加重复闹钟
        api.manager.AddRepeatingAlarm(req.Schedule, req.Message)
        c.JSON(http.StatusOK, AlarmResponse{
            Success: true,
            Message: "重复闹钟添加成功",
            Data: map[string]interface{}{
                "schedule": req.Schedule,
                "message":  req.Message,
            },
        })
    } else {
        c.JSON(http.StatusBadRequest, AlarmResponse{
            Success: false,
            Message: "必须提供时间或调度表达式",
        })
    }
}

// 处理移除闹钟的请求
func (api *AlarmAPI) handleRemoveAlarm(c *gin.Context) {
    idStr := c.Query("id")
    if idStr == "" {
        c.JSON(http.StatusBadRequest, AlarmResponse{
            Success: false,
            Message: "必须提供闹钟ID",
        })
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, AlarmResponse{
            Success: false,
            Message: "无效的ID格式: " + err.Error(),
        })
        return
    }

    api.manager.RemoveAlarm(cron.EntryID(id))
    c.JSON(http.StatusOK, AlarmResponse{
        Success: true,
        Message: "闹钟已移除",
    })
}

// 处理列出所有闹钟的请求
func (api *AlarmAPI) handleListAlarms(c *gin.Context) {
    alarms := make([]map[string]interface{}, 0)
    for _, alarm := range api.manager.alarms {
        alarmData := map[string]interface{}{
            "id":      alarm.ID,
            "time":    alarm.Time,
            "repeat":  alarm.Repeat,
            "message": alarm.Message,
        }
        alarms = append(alarms, alarmData)
    }

    c.JSON(http.StatusOK, AlarmResponse{
        Success: true,
        Message: "闹钟列表获取成功",
        Data:    alarms,
    })
}
