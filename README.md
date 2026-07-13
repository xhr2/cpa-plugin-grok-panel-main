# Grok Panel — CLIProxyAPI 原生插件

> 为 [CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) 打造的 Grok 账号管理面板，一个插件搞定统计、分类、健康检查和安全清理。

**版本** `v1.1.20` ｜ **平台** Linux / macOS / Windows ｜ **语言** 中文 ｜ **License** MIT

**仓库地址**：https://github.com/TizenryA/cpa-plugin-grok-panel

---

## 友链

> 感谢 [LINUX DO](https://linux.do/) 社区对开源项目的支持。本项目的开发与推广得益于社区环境，在此致敬。

[![LINUX DO](https://img.shields.io/badge/LINUX%20DO-社区友链-0066cc)](https://linux.do/)

---

## 这是什么

Grok Panel 是一个 CPA 原生 Go 插件（`.so` shared object），安装后直接在 CPA 管理中心的插件菜单中打开。它读取当前 CPA 实例中的所有 xAI / Grok auth 文件，提供可视化的账号管理面板。

**查看统计默认零配置** — 插件通过 CPA 官方 host callback 直接读取数据，每个人安装后只看到自己的账号。删除 / 检查若管理中心 iframe 读不到会话密钥，可在面板设置里本地保存管理密钥（仅当前浏览器）。

---

## 特点

- **即装即用**：通过 CPA 插件商店一键安装，无需手动下载或编译
- **默认零配置**：统计面板无需填写 CPA 地址；删除/检查可选用面板本地保存的管理密钥
- **安全隔离**：源码、Release 和配置中均不含任何密钥；不连接作者的 CPA；不上传账号数据
- **跟随系统主题**：圆润无衬线（Nunito / 系统 UI），暗色柔和蓝紫、亮色暖橙，圆角卡片自动跟随系统
- **移动端适配**：vw/vh 全比例布局，手机竖排单列
- **账号表分页**：默认每页 50 条，可选 20/50/100/200，号多也不卡
- **可分享**：插件仓库公开，其他 CPA 用户可直接安装使用

---

## 功能

### 统计概览

- Grok 文件总数、活跃 / 禁用账号数
- 总成功 / 失败请求数、成功率
- 估算 Token 用量与总容量（每账号默认 200 万 Token，可调）
- 使用率进度条（绿 / 黄 / 红三段）
- 请求趋势柱状图（成功绿色 / 失败红色）
- 账号类型分布（Free / Super / Heavy）
- 健康概览（正常 / 警告 / 无效 / 未知）

### 账号分类

从 auth 文件的 `note`、`prefix`、`subscription`、`plan` 等字段自动识别：

| 类型 | 说明 | 识别依据 |
|------|------|----------|
| **Free** | 普通免费账号 | 无 SuperGrok / Heavy 信号 |
| **Super** | SuperGrok 高级套餐 | `note: supergrok`、`prefix: supergrok` 等 |
| **Heavy** | Heavy 大用量套餐 | `note` / `prefix` 含 heavy 关键词 |

无法可靠判断时归为 Free（实际使用中绝大多数 CPA xAI 账号为免费账号）。

### 健康检查

- 点击"检查"对单个账号发起轻量探测
- "检查选中"批量检查勾选账号
- "手动检查本页"检查当前页账号

状态判定规则：

| 状态 | 条件 | 说明 |
|------|------|------|
| **正常** | 探测成功 | 绿色发光圆点 |
| **警告** | 429 限流 / 5xx / 超时 | 黄色圆点，不累计失效 |
| **无效** | 连续 3 次明确 401/403 | 红色圆点，可清理 |
| **禁用** | CPA 标记 disabled | 灰色圆点 |
| **未知** | 尚未检查 | 灰色空心圆点 |

> 429、5xx、网络超时**不会**被判定为失效。只有连续 3 次明确的认证失败（401/403）才成为自动删除候选。阈值可调。

### 删除与清理

- **单个删除**：点击"删除"→ 按钮变红显示"确认删除"→ 再次点击执行
- **批量删除**：勾选多个账号 → "删除选中"
- **一键清理**："清理无效 N"按钮，清理已确认失效且未受保护的账号
- **不弹窗**：所有确认通过按钮变色实现，不使用 `confirm` 弹窗

### 保护规则

默认配置（可调整）：

```
☑ 保护 Super    — SuperGrok 账号永不自动删除
☑ 保护 Heavy    — Heavy 账号永不自动删除
☑ 保护未知      — 未知类型账号永不自动删除
☐ 自动检查      — 默认关，开启后刷新时自动检查全部活跃账号
☐ 自动删除      — 默认关，开启后在自动检查后清理失效账号
```

受保护账号的删除按钮自动禁用，并在邮箱下方显示 ⚠ 标记。手动删除仍遵守保护开关。

### 筛选与搜索

- 搜索：邮箱、状态、类型、健康关键词
- 状态筛选：全部 / 活跃 / 禁用 / 其他 / 未知
- 类型筛选：全部 / Free / Super / Heavy
- 健康筛选：全部 / 健康 / 警告 / 无效 / 禁用 / 未知
- 用量筛选：未使用 / 低于一半 / 一半以上 / 高于八成
- 排序：成功请求 / 失败 / 用量 / 健康优先 / 类型 / 邮箱

### 设置

面板设置保存在当前浏览器本地：

- 每账号估算容量（默认 2,000,000 Token）
- 平均 Token / 请求（默认 5,000）
- 连续认证失败阈值（默认 3）
- 自动检查开关（默认关）
- 自动删除开关（默认关）
- Super / Heavy / 未知类型保护开关（默认开）

> 性能相关功能默认关闭，安装后不会立即批量检查账号。

---

## 安装方式

### 方式 A：插件商店安装（推荐）

**1. 配置 CPA**

在 CPA 的 `config.yaml` 中启用插件并添加本仓库为插件源：

```yaml
plugins:
  enabled: true
  dir: "plugins"
  store-sources:
    - "https://raw.githubusercontent.com/TizenryA/cpa-plugin-grok-panel/main/registry.json"
  configs:
    grok-panel:
      enabled: true
```

**2. 重启 CPA**

**3. 安装插件**

进入 CPA 管理中心 → 插件商店 → 找到 **Grok Panel** → 点击安装

或使用 API 安装：

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_MANAGEMENT_KEY" \
  "https://YOUR_CPA_HOST/v0/management/plugin-store/grok-panel/install?source=YOUR_SOURCE_ID&version=v1.1.20"
```

> 这里的管理密钥只用于执行安装操作，不会写入插件。

**4. 打开面板**

安装完成后，在 CPA 管理中心的插件菜单中点击 **Grok Panel** 即可打开。

### 方式 B：手动安装

安装包**只发布在** [GitHub Releases](https://github.com/TizenryA/cpa-plugin-grok-panel/releases)（仓库源码树不再存放 zip / `.so`）。下载与宿主匹配的压缩包：

```text
grok-panel_1.1.20_linux_amd64.zip
grok-panel_1.1.20_linux_arm64.zip
grok-panel_1.1.20_darwin_amd64.zip
grok-panel_1.1.20_darwin_arm64.zip
grok-panel_1.1.20_windows_amd64.zip
grok-panel_1.1.20_windows_arm64.zip
```

解压后将插件库文件放入 CPA 配置的插件目录：

```text
plugins/linux/amd64/grok-panel-v1.1.20.so
plugins/linux/arm64/grok-panel-v1.1.20.so
plugins/darwin/amd64/grok-panel-v1.1.20.dylib
plugins/darwin/arm64/grok-panel-v1.1.20.dylib
plugins/windows/amd64/grok-panel-v1.1.20.dll
plugins/windows/arm64/grok-panel-v1.1.20.dll
```

CPA 安装时会按宿主机 `GOOS/GOARCH` 自动选择对应 zip。当前已发布：

| 平台 | 资产 | 库文件 |
|---|---|---|
| Linux x86_64 | `*_linux_amd64.zip` | `grok-panel.so` |
| Linux ARM64 | `*_linux_arm64.zip` | `grok-panel.so` |
| macOS Intel | `*_darwin_amd64.zip` | `grok-panel.dylib` |
| macOS Apple Silicon | `*_darwin_arm64.zip` | `grok-panel.dylib` |
| Windows x86_64 | `*_windows_amd64.zip` | `grok-panel.dll` |
| Windows ARM64 | `*_windows_arm64.zip` | `grok-panel.dll` |

若出现 `plugin_install_failed`，优先检查：

1. 宿主架构是否已有对应 zip
2. CPA 机器是否能访问 GitHub Release
3. 返回体里的 `message` 原文（比 toast 更准确）

---

### 网络代理

插件对 Grok 官方订阅接口的出站请求会优先走 CPA 宿主回调 `host.http.do`，因此会自动使用 CPA `config.yaml` 中的 `proxy-url`（以及宿主侧代理策略）。无需在插件内单独配置代理。

若宿主不支持该回调，会回退到遵循系统环境变量 `HTTP_PROXY` / `HTTPS_PROXY` / `ALL_PROXY` 的本地 HTTP 客户端。

## 隐私与安全

- 不连接作者的 CPA，每个安装者只读取自己的数据
- 不上传账号文件或凭据
- 不在日志中打印 token、cookie 或授权头
- 删除使用 CPA 官方管理 API，不直接修改宿主文件系统
- 面板不向浏览器返回 access token、refresh token、SSO token

## 升级与卸载

- **升级**：在 CPA 插件商店点击更新，或重新调用安装接口指定新版本
- **卸载**：在 CPA 插件管理页面卸载 `grok-panel`，然后重启 CPA

## 构建

需要与 CPA 宿主兼容的 OS / CPU 架构及 Go 工具链（Go 1.24+）：

```bash
go test ./...
go vet ./...
# Linux 示例
go build -buildmode=c-shared -o grok-panel.so .
```

交叉编译产物与 zip **不要提交进仓库**，发版时上传到 GitHub Release 即可。仓库仅保留源码与 `registry.json`。

## GitHub Actions 自动编译

仓库已包含 `.github/workflows/build-release.yml`：

- **push 到 main/master** 或 **手动 Run workflow**：自动测试 → 编译 **linux/amd64** → **直接发布 GitHub Release**
- 版本号读取 `main.go` 里的 `pluginVersion`（例如 `1.1.21` → 标签 `v1.1.21`）
- 产物：`grok-panel_<version>_linux_amd64.zip` + `checksums.txt`

### 使用方式

```bash
# 确认 main.go 的 pluginVersion，以及 registry.json 的 version 一致
git add .
git commit -m "release: v1.1.21"
git push origin main
# push 成功后 Actions 自动编译并发布，无需再打 tag
```

CPA 的 `store-sources` 若指向你的仓库 `registry.json`，刷新插件商店即可安装新版本。

## 仓库结构

```text
main.go          # 插件后端
html.go          # 内嵌面板前端
main_test.go     # 单元测试
registry.json    # CPA 插件商店元数据
README.md
.gitignore
```

## License

MIT
