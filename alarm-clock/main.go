package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 创建闹钟管理器
	alarmManager := NewAlarmManager()
	if err := alarmManager.Start(); err != nil {
		log.Fatalf("闹钟管理器启动失败: %v", err)
	}

	// 创建API服务
	api := NewAlarmAPI(alarmManager)

	// 启动API服务器（非阻塞）
	go api.StartServer("8080")
	fmt.Println("闹钟服务已启动，API服务运行在 http://localhost:8080")

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 优雅关闭
	fmt.Println("正在关闭服务...")
	alarmManager.Stop()
}
