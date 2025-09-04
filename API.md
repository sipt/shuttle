# Shuttle API 文档

Shuttle 是一个功能完整的网络代理服务器，支持多种协议（HTTP/HTTPS/SOCKS）和灵活的路由规则配置。

## 概览

- **基础 URL**: `http://localhost:8080`
- **认证方式**: Basic Auth 或 Bearer Token
- **响应格式**: JSON
- **WebSocket 支持**: 实时数据推送

## 认证

### Basic Auth
```
Authorization: Basic <base64(username:password)>
```

### Bearer Token
```
Authorization: Bearer <token>
```

## 通用响应格式

### 成功响应
```json
{
  "code": 0,
  "message": "",
  "data": <响应数据>
}
```

### 错误响应
```json
{
  "code": 1,
  "message": "错误信息"
}
```

## API 接口

### 系统管理

#### 获取系统状态
```
GET /api/status
```

**响应示例:**
```json
{
  "code": 0,
  "data": "running"
}
```

**状态值:**
- `starting` - 启动中
- `running` - 运行中
- `stopped` - 已停止

#### 获取入站监听器列表
```
GET /api/inbounds
```

**响应示例:**
```json
{
  "code": 0,
  "data": [
    {
      "name": "http-proxy",
      "typ": "http",
      "addr": ":8080"
    }
  ]
}
```

#### 重新加载配置
```
PUT /api/reload
```

**响应示例:**
```json
{
  "code": 0,
  "data": "success"
}
```

### 服务器管理

#### 获取服务器列表
```
GET /api/servers
```

**响应示例:**
```json
{
  "code": 0,
  "data": [
    {
      "name": "proxy-server-1",
      "typ": "socks",
      "rtt": "100ms"
    }
  ]
}
```

#### 获取指定服务器信息
```
GET /api/servers/{name}
```

**参数:**
- `name` (path) - 服务器名称

**响应示例:**
```json
{
  "code": 0,
  "data": {
    "name": "proxy-server-1",
    "typ": "socks",
    "rtt": "100ms"
  }
}
```

#### 测试服务器延迟
```
PUT /api/servers/{name}/rtt
```

**参数:**
- `name` (path) - 服务器名称

### 服务器组管理

#### 获取服务器组列表
```
GET /api/groups
```

**响应示例:**
```json
{
  "code": 0,
  "data": [
    {
      "name": "auto-select",
      "typ": "select",
      "selected": {
        "name": "proxy-server-1",
        "typ": "socks",
        "rtt": "100ms",
        "selected": true
      },
      "servers": [
        {
          "name": "proxy-server-1",
          "typ": "socks",
          "rtt": "100ms",
          "selected": true
        }
      ]
    }
  ]
}
```

#### 获取指定服务器组信息
```
GET /api/group?group={name}
```

**参数:**
- `group` (query) - 服务器组名称

#### 选择服务器组中的服务器
```
PUT /api/group/select?group={group}&server={server}
```

**参数:**
- `group` (query) - 服务器组名称
- `server` (query) - 要选择的服务器名称

#### 重置服务器组延迟测试
```
PUT /api/group/rtt?group={group}
```

**参数:**
- `group` (query) - 服务器组名称

### DNS 管理

#### 获取 DNS 缓存列表
```
GET /api/dns/cache
```

**响应示例:**
```json
{
  "code": 0,
  "data": [
    {
      "typ": "dynamic",
      "domain": "example.com",
      "ip": "1.2.3.4",
      "dns_server": "udp://8.8.8.8:53",
      "country_code": "US"
    }
  ]
}
```

#### 清空 DNS 缓存
```
DELETE /api/dns/cache
```

### 规则管理

#### 获取当前路由模式
```
GET /api/rule/mode
```

**响应示例:**
```json
{
  "code": 0,
  "data": "ModeRule"
}
```

**路由模式:**
- `ModeDirect` - 直连模式
- `ModeGlobal` - 全局代理模式
- `ModeRule` - 规则代理模式

#### 设置路由模式
```
PUT /api/rule/mode/{mode}
```

**参数:**
- `mode` (path) - 路由模式

### 连接记录管理

#### 获取连接记录列表
```
GET /api/records
```

**响应示例:**
```json
{
  "code": 0,
  "data": [
    {
      "id": 12345,
      "dest_addr": "example.com:443",
      "policy": "PROXY",
      "up": 1024,
      "down": 2048,
      "status": "completed",
      "timestamp": 1640995200000,
      "protocol": "https",
      "duration": 5000,
      "dumped": true
    }
  ]
}
```

**连接状态:**
- `active` - 活跃连接
- `completed` - 连接完成
- `failed` - 连接失败
- `rejected` - 连接被拒绝

#### 清空连接记录
```
DELETE /api/records
```

### 连接管理

#### 获取当前连接列表
```
GET /api/conns
```

**响应示例:**
```json
[
  {
    "input": {
      "id": 12345,
      "local_addr": "127.0.0.1:8080",
      "remote_addr": "example.com:443"
    },
    "output": {
      "id": 12346,
      "local_addr": "192.168.1.100:12345",
      "remote_addr": "1.2.3.4:443"
    }
  }
]
```

### 流量转储管理

#### 获取流量转储状态
```
GET /api/stream/dump/status
```

**响应示例:**
```json
{
  "code": 0,
  "data": {
    "dump": true,
    "mitm": false
  }
}
```

#### 设置流量转储状态
```
PUT /api/stream/dump/status
```

**请求体:**
```json
{
  "dump": true,
  "mitm": true
}
```

#### 生成 CA 证书
```
POST /api/stream/dump/cert
```

#### 下载 CA 证书
```
GET /api/stream/dump/cert
```

**响应:** 二进制文件 (application/octet-stream)

#### 获取会话转储数据
```
GET /api/stream/dump/session/{id}
```

**参数:**
- `id` (path) - 连接记录 ID

**响应:** Base64 编码的 JSON 数据，包含请求和响应的完整信息

#### 获取请求体数据
```
GET /api/stream/dump/request/body/{id}
```

**参数:**
- `id` (path) - 连接记录 ID

**响应:** 原始请求体数据

#### 获取响应体数据
```
GET /api/stream/dump/response/body/{id}
```

**参数:**
- `id` (path) - 连接记录 ID

**响应:** 原始响应体数据

## WebSocket 接口

### RTT 延迟测试事件流
```
GET /ws/rtt
```

**消息格式:**
```json
{
  "name": "group_name",
  "typ": "select",
  "selected": {
    "name": "server_name",
    "typ": "socks",
    "rtt": "100ms"
  },
  "servers": [...]
}
```

### 连接记录事件流
```
GET /ws/records/events
```

**消息格式:**
```json
{
  "typ": "CreateRecordEvent",
  "value": {
    "id": 123,
    "dest_addr": "example.com:443",
    "policy": "PROXY",
    "up": 1024,
    "down": 2048,
    "status": "completed",
    "timestamp": 1640995200000,
    "protocol": "https",
    "duration": 5000,
    "dumped": true
  }
}
```

### 实时流量统计数据流
```
GET /ws/traffic
```

**消息格式:**
```json
{
  "up": 1024000,
  "down": 2048000
}
```

## 错误码说明

- `0` - 成功
- `1` - 一般错误（参数错误、资源不存在等）
- `500` - 内部服务器错误

## 使用示例

### 使用 curl 获取系统状态
```bash
curl -X GET "http://localhost:8080/api/status" \
  -H "Authorization: Basic <base64_credentials>"
```

### 使用 curl 设置路由模式
```bash
curl -X PUT "http://localhost:8080/api/rule/mode/ModeRule" \
  -H "Authorization: Bearer <token>"
```

### 使用 JavaScript 连接 WebSocket
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/traffic');
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Up:', data.up, 'Down:', data.down);
};
```

## 注意事项

1. 所有 API 接口都需要认证
2. WebSocket 连接也需要通过 HTTP 认证机制
3. 流量转储功能需要先生成 CA 证书
4. 连接记录的时间戳使用毫秒级 Unix 时间戳
5. RTT 延迟时间格式为字符串（如 "100ms"、"no rtt"、"failed"）


