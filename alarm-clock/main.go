package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // 创建闹钟管理器
    alarmManager := NewAlarmManager()
    alarmManager.Start()

    // 创建API服务
    api := NewAlarmAPI(alarmManager)

    // 启动API服务器（非阻塞）
    go api.StartServer("8080")

    // 等待中断信号以优雅地关闭服务器
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("关闭服务...")
}
