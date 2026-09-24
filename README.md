# RSS Notifier

一个用 Go 编写的轻量级 RSS/Atom 更新提醒程序。程序启动后立即检查订阅源，之后每 45 分钟检查一次；发现符合筛选条件的新条目时，通过桌面通知提醒，并将已处理条目记录在本地，避免重复通知。

本项目最初通过 Claude 网页版编写。

## 功能

- 支持 RSS、Atom 和 JSON Feed（由 `gofeed` 解析）。
- 可为每个订阅源设置标题关键词筛选。
- 使用 `beeep` 发送桌面通知。
- 将已见条目保存到 `~/.local/share/rss-notifier/seen.json`。
- 支持为订阅源配置本地图标，图标目录为 `~/.local/share/rss-notifier/icons/`。

## 配置订阅源

编辑 [`rsstest.go`](rsstest.go) 中的 `feeds` 列表。每项包含名称、Feed 地址、可选标题筛选词和可选图标文件名：

```go
{
    Name:   "Example Feed",
    URL:    "https://example.com/feed.xml",
    Filter: "release", // 留空表示不筛选
    Icon:   "example.png", // 留空表示使用系统默认图标
},
```

筛选不区分大小写。图标文件应放在 `~/.local/share/rss-notifier/icons/` 下。

## 运行

需要安装 Go（版本以 `go.mod` 为准）以及当前桌面环境可用的通知服务。然后在项目目录执行：

```sh
go run .
```

也可以编译后运行：

```sh
go build -o rss-notifier .
./rss-notifier
```

程序会在启动时立即检查一次，随后每 45 分钟检查。运行期间可用 `Ctrl+C` 退出。

## 本地依赖说明

当前 [`go.mod`](go.mod) 将 `github.com/mmcdole/gofeed` 替换为 `../gofeed`。因此构建时需要在本项目的同级目录提供对应的 `gofeed` 源码目录；若要独立分发项目，请将该替换改为正式模块版本或调整为适合分发的依赖方式。

## 数据与图标位置

- 已见条目：`~/.local/share/rss-notifier/seen.json`
- 图标：`~/.local/share/rss-notifier/icons/`

Feed 列表和检查间隔目前直接在源码中配置。当前默认订阅 Linux 内核 LTS 及 Arch Linux 新闻源。
