# omnigram-server

omnigram-server 是使用 Golang 编写的个人内容管理后端服务，提供类似 jellyfin 和 navidrome 文档搜刮与管理功能。
客户端支持请前往 [Omnigram](github.com/nexptr/omnigram).

> 提醒：中国境内网站，个人是不允许进行在线出版的，维护公开的书籍网站是违法违规的行为！建议仅作为个人使用！

### 目标与功能介绍

1. 提供私有化电子书库搜刮服务（当前仅支持Epub格式）；
2. 使用TTS大模型实现，提供文字转语音服务（不同模型需要的硬件配置参考 [](tts.md)）；
3. 提供Chat模型支持对电子书进行理解分析，支持对自己简单日记进行扩展书写和润色；

## 部署方式

### 快速启动

### docker

TODO

### linux server

使用Mysql或者Postgresql数据库

1. 修改配置文件 conf.yaml

### 开发计划

- [] 支持音频、视频搜刮；
- [] 支持浏览管理
- [] 支持多用户管理与用户分级

## 常见问题

常见问题请参阅使用指南

手动安装请参考[开发指南](docs/dev.md)

## 项目依赖以及三分框架使用说明

TODO

## 鸣谢
