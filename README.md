# search

`bamgoo` 全文搜索模块，支持统一 API + 多驱动。

## 能力

1. 多连接配置 + 权重分布。
2. 统一索引接口：
   - 索引按 `Register(indexName, search.Index{...})` 自动同步
   - `Upsert(index, rows ...Map)`
   - `Clear(index)`
   - `Delete`
3. 统一查询接口：
   - `Search(index, keyword, args ...Any)`
   - `Count(index, keyword, args ...Any)`
   - `Signature(index, keyword, args ...Any)`
4. 统一查询 DSL（Map，推荐 `$` 前缀）：
   - `$filters`
   - `$sort`
   - `$offset/$limit`
   - `$fields`（同 `$select`）
   - `$facets`
   - `$highlight`
   - `$prefix`（前缀匹配）
   - 高亮内容直接写回 `hits[].payload` 对应字段
5. 支持索引字段定义（`Index.Attributes`）：
   - 写入前按字段定义做 `Mapping`
   - 查询结果按字段定义做 `Mapping`
   - 驱动建索引时可按属性生成结构（驱动支持范围内）
6. 可查看驱动能力：
   - `GetCapabilities(index)`
   - `ListCapabilities()`

## 注册与配置

```go
import (
    "github.com/bamgoo/bamgoo"
    "github.com/bamgoo/search"
)

bamgoo.Register("default", search.Config{
    Driver: "file",
    Weight: 1,
    Prefix: "demo",
})

bamgoo.Register("article", search.Index{
    Primary: "id",
    Attributes: Vars{
        "id":       Var{Type: "string", Required: true},
        "title":    Var{Type: "string", Required: true},
        "content":  Var{Type: "string"},
        "category": Var{Type: "string"},
        "score":    Var{Type: "float"},
        "created":  Var{Type: "timestamp"},
    },
})
```

也可通过配置文件：

```toml
[search]
driver = "file"
weight = 1
prefix = "demo"

[search.setting]
path = "data/search"
```

## 查询示例

```go
res, err := search.Search("article", "go", Map{
    "$filters": Map{
        "category": "tech",
        "score": Map{"$gte": 8.5},
    },
    "$sort": Map{"score": DESC},
    "$offset": 0,
    "$limit":  20,
    "$fields": []string{"title", "category", "score"},
    "$facets": []string{"category"},
    "$highlight": []string{"title", "content"},
    "$prefix": false,
})
_ = res
_ = err
```

## filter 约定

支持以下操作符（统一使用 `$` 前缀）：

- `$eq` / `$ne`
- `$in` / `$nin`
- `$gt` / `$gte` / `$lt` / `$lte`
- `$range`（配合 `min/max`）

只支持 Map 风格参数，旧格式已移除：

- sort 的 `field/desc`
- filter 的 `op/value`

示例：

```go
Map{"$filters": Map{"category": "123"}}
Map{"$filters": Map{"category": Map{"$eq": 123}}}
Map{"$filters": Map{"score": Map{"$gt": 100, "$lt": 500}}}
Map{"$sort": Map{"score": DESC}}
Map{"$sort": []Map{{"score": DESC}, {"id": ASC}}}
Map{"category": "tech", "$sort": Map{"score": DESC}} // 顶层字段也可直接当过滤条件
```

## 驱动

- `search-file`（driver: `file`）
- `search-meilisearch`（driver: `meilisearch`/`meili`）
- `search-opensearch`（driver: `opensearch`）
- `search-elasticsearch`（driver: `elasticsearch`/`es`）
