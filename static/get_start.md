## 快速开始
### 简单配置
修改shuttle.yaml中的`Proxy`和`Proxy-Group`配置项，改成自己的服务器。
示例配置中的`Rule`中用到了分组名`Proxy`,`nProxy`，所以要保留。
```
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
```
### 安装与启动
#### Mac & Linux
1. 启动方法
```
cd shuttle
./start.sh #不会有任何输出内容
```
2. 根据配置文件，配置系统网络代理：HTTP Proxy、HTTPS Proxy、Socks Proxy
3. 设置完成后可以访问`http://c.sipt.top`查看到控件台
4. 关闭方法：web_ui里点击shutdown。

#### Windows
1. 启动方法：双击`startup.bat`
2. 根据配置文件，配置系统网络代理
3. 设置完成后可以访问`http://c.sipt.top`查看到控件台
4. 关闭方法：web_ui里点击shutdown。