# RSS Notifier

[English](README.md) | **简体中文**

一个用 Go 编写的轻量级 RSS/Atom 更新提醒程序。程序启动后立即检查订阅源，之后按配置的间隔检查；符合筛选条件的新文章会触发桌面通知，并记录已处理条目以避免重复通知。通知中会包含文章链接。
在 Windows 和 Linux 上，点击通知会使用默认浏览器打开文章；点击“Copy link”即可将文章链接复制到剪贴板。

本项目最初由 Claude 网页版协助编写。

## 功能

- 使用 `gofeed` 解析 RSS、Atom 和 JSON Feed
- 可为每个订阅源设置不区分大小写的标题筛选词
- 在 Windows 和 Linux 上发送桌面通知，并提供复制文章链接的操作
- 可配置点击通知时是否在默认浏览器中打开文章
- 可配置通知停留时间；支持的平台可让通知保持显示，直到手动关闭
- 可配置复制链接后是否关闭通知
- 在 Linux 上可配置通知紧急级别（`low`、`normal` 或 `critical`）；具体样式由桌面环境决定
- 将已见条目保存在 `~/.local/share/rss-notifier/seen.json`
- 支持本地通知图标，图标目录为 `~/.local/share/rss-notifier/icons/`

## 订阅配置

配置文件位于 `~/.local/share/rss-notifier/feeds.json`。首次运行时，程序会生成默认配置，检查间隔为 45 分钟，订阅列表为空。编辑该文件添加订阅或调整检查间隔：

```json
{
  "check_interval": "45m",
  "notification_timeout": "5s",
  "notification_urgency": "normal",
  "open_link_on_click": true,
  "dismiss_on_copy": true,
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

- 选项说明
   - `check_interval`：Go 时长格式（如 `60s`、`30m` 或 `1h`），必须大于零
   - `notification_timeout`：接受非负 Go 时长（如 `5s`）或 `"never"`（通知将保持显示直至手动关闭），默认值为 `5s`
     - Linux 会直接应用指定时长，但是当 `notification_urgency` 为 `critical` 时，会按通知规范不自动过期，此时该选项失效；
     - Windows 系统通知仅支持“短”和“长”两档，有限时长会被映射为系统预设时长，`"never"` 则使用 Windows 的提醒（reminder）模式。此外，程序向后兼容旧版纯数组订阅格式（默认检查间隔为 45 分钟，通知停留时间使用默认值）；省略 `filter` 或设为空字符串表示不过滤标题；省略 `icon` 则使用系统默认图标，自定义图标请放置于 `~/.local/share/rss-notifier/icons/`
   - `notification_urgency`：在 Linux 上可设为 `low`、`normal` 或 `critical`。视觉效果取决于桌面环境（例如 critical 可能使用颜色强调，或保持常驻直至手动关闭）；Windows 会忽略此设置
   - `open_link_on_click`：控制点击通知时是否在默认浏览器中打开文章，默认值为 `true`。设为 `false` 可禁用该行为，同时仍保留“Copy Link”操作
   - `dismiss_on_copy`：控制复制链接后是否关闭通知，默认值为 `true`；设为 `false` 可让通知继续显示。注意部分 Linux 通知服务可能不支持 `resident` 提示；Windows Toast 的按钮操作通常会自动关闭通知，程序会在复制后重新显示该通知，因此不建议在 Windows 上设为 `false`

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

程序运行时用到的文件集中保存在同一个数据目录中：

```text
~/.local/share/rss-notifier/
├── feeds.json       # 订阅源和通知设置
├── seen.json        # 已处理的文章，用于避免重复提醒
├── initialized      # 标记首次静默扫描已完成
└── icons/           # feeds.json 中引用的本地通知图标
```

在 Windows 上，`~` 代表当前用户目录，例如 `C:\Users\Your Name`。若要使用图标，在订阅项的 `icon` 字段中填写文件名，并将图标放入 `icons/`；例如 `icon: "example.png"` 对应 `~/.local/share/rss-notifier/icons/example.png`。
