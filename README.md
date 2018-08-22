# Shuttle



![Shuttle](./Shuttle_Logo.PNG)

感谢logo提供者：**@不二**

有疑问可以 @sipt(wxysipt@gmail.com)，也可以`Issues`

代码整理中，即将开源...

## 介绍
**Shuttle**是一个基于Go开发的**全平台**ss-local工具，具有代理、多服务器选择、HTTP/HTTPS抓包、独立DNS解析机制，目标为开发者提供便利。

参照软件**Surge for Mac**。

![Introduction](static/dump_mitm.jpg)

[快速开始](static/get_start.md)

[配置文件](static/config.md)

[简陋的web-ui](static/web_ui.md)


## 功能
- [ ] 代理功能
	- [x] TCP(HTTP/HTTPS)
	- [ ] UDP
- [x] 扩展功能
	- [x] HTTP抓包
	- [x] HTTPS抓包(MITM)
	- [x] keep-alive时请求切分
	- [ ] 请求头修改
	- [ ] 返回头修改
	- [ ] 请求mapping
- [x] 远端多服务器管理
	- [x] 服务器分组包含
	- [x] 服务器选择
		- [x] RTT(往返时间)选择
		- [x] Select(手动)选择
- [ ] 代理模式
	- [ ] 全局代理
	- [x] 规则代理
		- [x] DOMAIN：域名全匹配
		- [x]  DOMAIN-SUFFIX：域名后缀匹配
		- [x]  DOMAIN-KEYWORD：域名关键字匹配
		- [x]  IP-CIDR：ip段匹配
		- [x]  GEO-IP: 支持GEO-IP路由
		- [ ]  USER-AGENT：HTTP头字匹配
- [x] DNS
	- [x] static：静态地址映射
	- [x] direct：直连DNS解析
	- [x] remote：远程服务器DNS解析(防止DNS污染)
	- [x] GEO-IP判断
- [x] 外部窗口
	- [x] API
		- [x]  获取服务器列表
		- [x]  RTT分组刷新
		- [x]  Select分组手动选择
		- [x]  DNS缓存获取
		- [x]  DNS缓存刷新
		- [x]  请求记录列表获取
		- [x]  请求记录清空
		- [x]  CA证书生成
		- [x]  CA证书下载
		- [x]  HTTP Dump开关
		- [x]  MITM 开关
		- [x]  HTTP/HTTPS抓包内容获取 
		- [x]  关闭Shuttle
		- [x]  重载配置
		- [ ]  全局代理开关
		- [ ]  支持Websocket，完成内容增量更新
	- [x] Web UI
		- [x] 很简陋的Web UI (angular6 + ant design)
- [ ] 优化
	- [ ] 内存优化
	- [ ] log日志


