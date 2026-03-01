# search

`search` 是 infrago 的**模块**。

## 包定位

- 类型：模块
- 作用：搜索模块，负责统一检索接口与多搜索后端适配。

## 主要功能

- 对上提供统一模块接口
- 对下通过驱动接口接入具体后端
- 支持按配置切换驱动实现

## 快速接入

```go
import _ "github.com/infrago/search"
```

```toml
[search]
driver = "default"
```

## 驱动实现接口列表

以下接口由驱动实现（来自模块 `driver.go`）：

### Driver

- `Connect(*Instance) (Connection, error)`

### Connection

- `Open() error`
- `Close() error`
- `Capabilities() Capabilities`
- `SyncIndex(name string, index Index) error`
- `Clear(index string) error`
- `Upsert(index string, rows []Map) error`
- `Delete(index string, ids []string) error`
- `Search(index string, query Query) (Result, error)`
- `Count(index string, query Query) (int64, error)`

## 全局配置项（所有配置键）

配置段：`[search]`

- 未检测到配置键（请查看模块源码的 configure 逻辑）

## 说明

- `setting` 一般用于向具体驱动透传专用参数
- 多实例配置请参考模块源码中的 Config/configure 处理逻辑
