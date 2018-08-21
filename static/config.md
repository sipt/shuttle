## 配置
```yaml
ver: v1.0.0 # 配置文件版本
General: #通用配置项
  loglevel: "info" # trace, debug, info, error;建议别开trace
  dns-server: # DNS服务器
  - "114.114.114.114"
  - "223.5.5.5"
  http-port: "8080" # httpProxy监听端口
  http-interface: "0.0.0.0" # 允许访问
  socks-port: "8081"
  socks-interface: "0.0.0.0"
  controller-port: "8082" # api/web ui端口
  controller-interface: "0.0.0.0"
Proxy: #服务器配置
  # 服务器名：[服务器地址域名/ip, 端口, 加密方式, 密码]
  "🇯🇵jp_a": ["jp.a.example.com", "12345", "rc4-md5", "123456"]
  "🇯🇵jp_b": ["jp.b.example.com", "12345", "rc4-md5", "123456"]
  "🇯🇵jp_c": ["jp.c.example.com", "12345", "rc4-md5", "123456"]
  "🇭🇰HK_b": ["hk.a.example.com", "12345", "rc4-md5", "123456"]
  "🇭🇰HK_b": ["hk.b.example.com", "12345", "rc4-md5", "123456"]
  "🇭🇰HK_c": ["hk.c.example.com", "12345", "rc4-md5", "123456"]
  "🇺🇸US_a": ["us.a.example.com", "12345", "rc4-md5", "123456"]
  "🇺🇸US_b": ["us.b.example.com", "12345", "rc4-md5", "123456"]
  "🇺🇸US_c": ["hk.c.example.com", "12345", "rc4-md5", "123456"]
Proxy-Group: #服务器分组配置
  ### 组名: [选择方式, 服务器/分组名 ...]
  "Auto": ["rtt", "🇭🇰HK_a", "🇭🇰HK_b", "🇭🇰HK_c",
  "🇯🇵JP_a", "🇯🇵JP_b", "🇯🇵JP_c",
  "🇺🇸US_a", "🇺🇸US_b", "🇺🇸US_c"]
  "HK": ["select", "🇭🇰HK_a", "🇭🇰HK_b", "🇭🇰HK_c"]
  "JP": ["select", "🇯🇵JP_a", "🇯🇵JP_b", "🇯🇵JP_c"]
  "US": ["select", "🇺🇸US_a", "🇺🇸US_b", "🇺🇸US_c"]
  "Proxy": ["select", "Auto", "US", "HK", "JP"]
  "nProxy": ["select", "DIRECT"]
Local-DNS: # DNS配置
# - [匹配方式，域名，解析方式，解析方式对应值]
# - [域名全匹配，域名，static(静态解析)，直接对应IP]
- ["DOMAIN", "localhost", "static", "127.0.0.1"]
# - [域名关键字匹配，关键字，remote(远程解析)，无]
- ["DOMAIN-KEYWORD", "google", "remote", ""]
# - [域名后缀匹配，后缀，direct(直连DNS服务器解析)，DNS服务器地址]
- ["DOMAIN-SUFFIX", "appspot.com", "direct", "114.114.114.114"]
MITM: # CA证书和私钥，不需要配置，由程序自动生成，保存在这里
  ca: (base64)
  key: (base64)
Rule: # 代理规则
# - [匹配方式，域名，连接方式，备注]
# - [域名后缀匹配，后缀，直连，]
- ["DOMAIN-SUFFIX", "gitlab.anjian.com", "DIRECT", ""]
# - [域名全匹配，域名，走分组Proxy，]
- ["DOMAIN", "sipt.top", "Proxy", ""]
# - [域名关键字匹配，关键字，拒绝连接，]
- ["DOMAIN-KEYWORD", "zjtoolbar", "REJECT", ""]
# - [IP网段断匹配，IP网段，直连，]
- ["IP-CIDR", "127.0.0.0/8", "DIRECT", ""]
# - [GEOIP匹配，中国，走nProxy组规则，]
- ["GEOIP", "CN", "nProxy", ""]
# - [以上都不满足，，走Proxy组规则，]
- ["FINAL", "", "Proxy", ""]
```
在realse版本中已经加入了`example.yaml`配置可供参考。
1. 加密方式支持：
 - aes-128-cfb
 - aes-192-cfb
 - aes-256-cfb
 - aes-128-ctr
 - aes-192-ctr
 - aes-256-ctr
 - des-cfb
 - bf-cfb
 - cast5-cfb
 - rc4-md5
 - chacha20
 - chacha20-ietf
 - salsa20
2. 选择方式：
 - select：手动选择
 - rtt：本机穿过远端到达`www.gstatic.com`的往返时间评出最优
