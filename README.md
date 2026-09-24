# RSS Notifier

一个用 Go 编写的轻量级 RSS/Atom 更新提醒程序。程序启动后立即检查订阅源，之后按配置的间隔检查；发现符合筛选条件的新条目时，通过桌面通知提醒，并将已处理条目记录在本地，避免重复通知。通知内容包含文章链接。

本项目最初通过 Claude 网页版编写。

## 功能

- 支持 RSS、Atom 和 JSON Feed（由 `gofeed` 解析）。
- 可为每个订阅源设置标题关键词筛选。
- 使用 `beeep` 发送桌面通知。
- 将已见条目保存到 `~/.local/share/rss-notifier/seen.json`。
- 支持为订阅源配置本地图标，图标目录为 `~/.local/share/rss-notifier/icons/`。

## 配置订阅源

订阅配置位于 `~/.local/share/rss-notifier/feeds.json`。首次运行时，程序会在该目录生成默认配置（45 分钟间隔、空订阅列表）；编辑该文件添加订阅或调整间隔。配置包含检查间隔和订阅列表；每项订阅包含名称、Feed 地址、可选标题筛选词和可选图标文件名：

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

`check_interval` 使用 Go 时长格式，例如 `30m` 或 `1h`，且必须大于零。旧版纯数组配置仍可读取，默认间隔为 45 分钟。`filter` 留空或省略表示不筛选，筛选不区分大小写。`icon` 留空或省略表示使用系统默认图标。图标文件应放在 `~/.local/share/rss-notifier/icons/` 下。

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

首次启动会静默记录订阅源当前已有的文章作为初始状态，后续新文章才发送通知；之后按 `check_interval` 指定的间隔检查。运行期间可用 `Ctrl+C` 退出。

## 数据与图标位置

- 订阅配置：`~/.local/share/rss-notifier/feeds.json`
- 已见条目：`~/.local/share/rss-notifier/seen.json`
- 图标：`~/.local/share/rss-notifier/icons/`

首次创建配置时，检查间隔默认为 45 分钟，可在 `feeds.json` 的 `check_interval` 中修改。
