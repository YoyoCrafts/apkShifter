# APKShifter

**APKShifter** 是一款专为防止 APK 报毒设计的工具，通过自动更换 APK 文件的包名与签名来避免被杀毒软件误报。该工具使用 Golang 开发，具备强大的处理能力，同时提供了简单易用的部署流程，帮助您轻松应对 APK 修改的需求。

## 功能特点

- **自动更换包名**：动态生成并替换 APK 文件的包名，防止应用被重复识别。
- **自动更换签名**：自动生成并替换 APK 文件的签名，无需手动干预，确保应用安全性。
- **高性能处理**：采用 Golang 开发，利用其高效并发处理能力，快速完成 APK 的修改任务。
- **简单部署**：轻量级架构，支持快速部署，无需复杂配置或依赖环境，开箱即用。

## 安装与使用
```bash
bash <(curl -sSL https://github.com/YoyoCrafts/apkShifter/edit/main/install.sh) 
```
- 安装成功后 
- 通过 "http://ip:port/guide/download/任意名称.apk" 下载出来过的app


- 如果apk已经集成了walle的需要自动打渠道包的话  
- 通过"http://ip:port/channel/download/渠道信息/任意名称.apk" 将自动写入渠道信息

### 环境要求

- 操作系统：Linux 服务器
- 最低配置：2 核心 CPU，2GB 内存

