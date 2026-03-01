# search

`search` 是 infrago 的模块包。

## 安装

```bash
go get github.com/infrago/search@latest
```

## 最小接入

```go
package main

import (
    _ "github.com/infrago/search"
    "github.com/infrago/infra"
)

func main() {
    infra.Run()
}
```

## 配置示例

```toml
[search]
driver = "default"
```

## 公开 API（摘自源码）

- `func BuildQuery(keyword string, args ...Any) Query`
- `func FilterMatch(filter Filter, payload Map) bool`
- `func QuerySignature(index string, q Query) string`
- `func (d *defaultDriver) Connect(inst *Instance) (Connection, error)`
- `func (c *defaultConnection) Open() error  { return nil }`
- `func (c *defaultConnection) Close() error { return nil }`
- `func (c *defaultConnection) Capabilities() Capabilities`
- `func (c *defaultConnection) SyncIndex(name string, index Index) error`
- `func (c *defaultConnection) Clear(name string) error`
- `func (c *defaultConnection) Upsert(index string, rows []Map) error`
- `func (c *defaultConnection) Delete(index string, ids []string) error`
- `func (c *defaultConnection) Search(index string, query Query) (Result, error)`
- `func (c *defaultConnection) Count(index string, query Query) (int64, error)`
- `func RegisterDriver(name string, driver Driver)`
- `func RegisterConfig(name string, cfg Config)`
- `func RegisterConfigs(configs Configs)`
- `func RegisterIndex(name string, index Index)`
- `func RegisterIndexes(indexes Indexes)`
- `func Clear(index string) error`
- `func GetCapabilities(index string) Capabilities`
- `func ListCapabilities() map[string]Capabilities`
- `func Upsert(index string, rows ...Map) error`
- `func Delete(index string, ids []string) error`
- `func Search(index, keyword string, args ...Any) (Result, error)`
- `func Count(index, keyword string, args ...Any) (int64, error)`
- `func Signature(index, keyword string, args ...Any) string`
- `func (m *Module) Register(name string, value Any)`
- `func (m *Module) RegisterDriver(name string, driver Driver)`
- `func (m *Module) RegisterConfig(name string, cfg Config)`
- `func (m *Module) RegisterConfigs(configs Configs)`
- `func (m *Module) RegisterIndex(name string, index Index)`
- `func (m *Module) RegisterIndexes(indexes Indexes)`
- `func (m *Module) Config(global Map)`
- `func (m *Module) Setup() {}`
- `func (m *Module) Open()`
- `func (m *Module) Start()`
- `func (m *Module) Stop() {}`
- `func (m *Module) Close()`
- `func (m *Module) Clear(index string) error`
- `func (m *Module) Capabilities(index string) Capabilities`

## 排错

- 模块未运行：确认空导入已存在
- driver 无效：确认驱动包已引入
- 配置不生效：检查配置段名是否为 `[search]`
