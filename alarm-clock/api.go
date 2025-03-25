package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

// API 结构体处理HTTP请求
type AlarmAPI struct {
	manager *AlarmManager
	router  *gin.Engine
}

// NewAlarmAPI 创建新的API处理器
func NewAlarmAPI(manager *AlarmManager) *AlarmAPI {
	router := gin.Default()
	api := &AlarmAPI{
		manager: manager,
		router:  router,
	}
	api.setupRoutes()
	return api
}

// 设置路由
func (api *AlarmAPI) setupRoutes() {
	alarmGroup := api.router.Group("/alarm")
	{
		alarmGroup.POST("/add", api.handleAddAlarm)
		alarmGroup.DELETE("/remove", api.handleRemoveAlarm)
		alarmGroup.GET("/list", api.handleListAlarms)
	}
}

// StartServer 启动HTTP服务器
func (api *AlarmAPI) StartServer(port string) error {
	return api.router.Run(":" + port)
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

	var err error
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

		err = api.manager.AddOneTimeAlarm(t, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AlarmResponse{
				Success: false,
				Message: "添加闹钟失败: " + err.Error(),
			})
			return
		}

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
		err = api.manager.AddRepeatingAlarm(req.Schedule, req.Message)
		if err != nil {
			c.JSON(http.StatusInternalServerError, AlarmResponse{
				Success: false,
				Message: "添加闹钟失败: " + err.Error(),
			})
			return
		}

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

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, AlarmResponse{
			Success: false,
			Message: "无效的ID格式: " + err.Error(),
		})
		return
	}

	if err := api.manager.RemoveAlarm(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, AlarmResponse{
			Success: false,
			Message: "删除闹钟失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AlarmResponse{
		Success: true,
		Message: "闹钟已移除",
	})
}

// 处理列出所有闹钟的请求
func (api *AlarmAPI) handleListAlarms(c *gin.Context) {
	alarms, err := api.manager.ListAlarms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, AlarmResponse{
			Success: false,
			Message: "获取闹钟列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, AlarmResponse{
		Success: true,
		Message: "闹钟列表获取成功",
		Data:    alarms,
	})
}
