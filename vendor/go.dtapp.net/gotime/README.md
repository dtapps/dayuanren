<h1>
<a href="https://www.dtapp.net/">Golang</a>
</h1>

📦 Golang 时间

[comment]: <> (go)
[![godoc](https://pkg.go.dev/badge/go.dtapp.net/gotime?status.svg)](https://pkg.go.dev/go.dtapp.net/gotime)
[![goproxy.cn](https://goproxy.cn/stats/go.dtapp.net/gotime/badges/download-count.svg)](https://goproxy.cn/stats/go.dtapp.net/gotime)
[![goreportcard.com](https://goreportcard.com/badge/go.dtapp.net/gotime)](https://goreportcard.com/report/go.dtapp.net/gotime)
[![deps.dev](https://img.shields.io/badge/deps-go-red.svg)](https://deps.dev/go/go.dtapp.net%2Fgotime)

#### 安装

```shell
go get -v -u go.dtapp.net/gotime@v1.0.9
```

#### 使用

```go
package main

import (
	"go.dtapp.net/gotime"
	"testing"
)

// TestVerification 验证字符串是否为时间
func TestVerification(t *testing.T) {
	t.Log(gotime.Verification("2022-02-05 00:00:00", gotime.DateTimeFormat))
}
```