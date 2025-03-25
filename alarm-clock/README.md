# 闹钟服务 API

这是一个基于 Go 语言开发的闹钟服务，使用 Gin 框架提供 HTTP API 接口供外部服务调用。

## 功能

- 添加一次性闹钟提醒
- 添加基于 cron 表达式的重复性闹钟提醒
- 移除已有闹钟
- 列出所有闹钟

## 依赖库

- [Gin](https://github.com/gin-gonic/gin) - Web 框架
- [Cron](https://github.com/robfig/cron) - 定时任务库

## 安装与运行

```bash
# 初始化Go模块
go mod init alarm-clock

# 安装依赖
go get github.com/gin-gonic/gin github.com/robfig/cron/v3

# 运行服务
go run .
```

服务将在本地 8080 端口启动。

## API 接口说明

### 1. 添加闹钟

**请求：**

- 方法：POST
- 路径：/alarm/add
- 内容类型：application/json

**请求体格式：**

添加一次性闹钟：

```json
{
  "time": "2023-05-20 08:30:00",  // 格式：yyyy-MM-dd HH:mm:ss
  "message": "起床提醒"
}
```

添加重复性闹钟：

```json
{
  "schedule": "0 30 8 * * *",  // cron表达式：秒 分 时 日 月 星期
  "message": "每天早上8点30分提醒"
}
```

**成功响应示例：**

```json
{
  "success": true,
  "message": "闹钟添加成功",
  "data": {
    "time": "2023-05-20T08:30:00Z",
    "message": "起床提醒"
  }
}
```

### 2. 移除闹钟

**请求：**

- 方法：DELETE
- 路径：/alarm/remove?id=1

**成功响应示例：**

```json
{
  "success": true,
  "message": "闹钟已移除"
}
```

### 3. 列出所有闹钟

**请求：**

- 方法：GET
- 路径：/alarm/list

**成功响应示例：**

```json
{
  "success": true,
  "message": "闹钟列表获取成功",
  "data": [
    {
      "id": 1,
      "time": "2023-05-20T08:30:00Z",
      "repeat": "none",
      "message": "起床提醒"
    },
    {
      "id": 2,
      "time": "0001-01-01T00:00:00Z",
      "repeat": "0 30 8 * * *",
      "message": "每天早上8点30分提醒"
    }
  ]
}
```

## Cron 表达式说明

Cron 表达式格式为 `秒 分 时 日 月 星期`，例如：

- `0 0 8 * * *` - 每天早上8点整
- `0 30 9 * * 1-5` - 工作日（周一至周五）早上9点30分
- `0 0 12 1 * *` - 每月1号中午12点

详细语法请参考 [cron 表达式文档](https://github.com/robfig/cron)。
