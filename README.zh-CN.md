# RSS Notifier

[English](README.md) | **简体中文**

一个用 Go 编写的轻量级 RSS/Atom 更新提醒程序。程序启动后立即检查订阅源，之后按配置的间隔检查；符合筛选条件的新文章会触发桌面通知，并记录已处理条目以避免重复通知。通知中会包含文章链接。

本项目最初由 Claude 网页版协助编写。

## 功能

- 使用 `gofeed` 解析 RSS、Atom 和 JSON Feed。
- 可为每个订阅源设置不区分大小写的标题筛选词。
- 使用 `beeep` 发送桌面通知。
- 将已见条目保存在 `~/.local/share/rss-notifier/seen.json`。
- 支持本地通知图标，图标目录为 `~/.local/share/rss-notifier/icons/`。

## 订阅配置

配置文件位于 `~/.local/share/rss-notifier/feeds.json`。首次运行时，程序会生成默认配置，检查间隔为 45 分钟，订阅列表为空。编辑该文件添加订阅或调整检查间隔：

```json
{
  "check_interval": "45m",
  "feeds": [
    {
      "name": "Example Feed",
      "url": "https://example.com/feed.xml",
      "filter": "release",
      "icon": "example.png"
    }
  ]
}
```

`check_interval` 使用 Go 时长格式，例如 `30m` 或 `1h`，且必须大于零。程序仍兼容旧版纯数组订阅配置，此时检查间隔默认为 45 分钟。省略 `filter` 或设为空字符串表示不过滤标题。省略 `icon` 则使用系统默认图标。图标文件请放在 `~/.local/share/rss-notifier/icons/` 下。

## 运行

安装 Go（所需版本以 `go.mod` 为准），并确保桌面通知服务可用。在项目目录执行：

```sh
go run .
```

也可以编译后运行：

```sh
go build -o rss-notifier .
./rss-notifier
```

首次运行时，程序会静默记录订阅源已有文章作为初始状态；初始化完成后，新文章才会触发通知。之后程序会按 `check_interval` 指定的间隔检查。按 `Ctrl+C` 退出。

## 数据与图标

- 订阅配置：`~/.local/share/rss-notifier/feeds.json`
- 已见条目：`~/.local/share/rss-notifier/seen.json`
- 图标：`~/.local/share/rss-notifier/icons/`
