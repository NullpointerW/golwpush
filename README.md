# golwpush
[![Go
Version](https://img.shields.io/github/go-mod/go-version/rwxrob/structs)](https://tip.golang.org/doc/go1.18)
[![Go Report
Card](https://goreportcard.com/badge/github.com/NullpointerW/golwpush)](https://goreportcard.com/report/github.com/NullpointerW/golwpush)
## 特性
 * 轻量级
 * 高性能
 * 纯Golang实现
 * 消息确认机制
 * 支持丢失消息持久化/重传
 * 支持单个/多个/广播推送
 * 通过聚合广播消息合并发送，减少io调用，大幅提升网络吞吐量
 * 心跳支持
 * 客户端连接的统计信息
 
 ## 安装

golang1.18，基于gomod管理依赖。

* 下载golwpush

```
go get github.com/NullpointerW/golwpush
```

* 安装依赖

```
export GOPROXY=goproxy.io
go mod download
```

* 运行service服务

```
 go run ./app 
```

* 运行client客户端
```
 go run  ./app/clinet &&
```
