> **语言**：本文档为简体中文翻译，[点此查看英文原版 README.md](README.md)。

[![Gitbook](https://badges.aleen42.com/src/gitbook_1.svg)](https://daseldocs.tomwright.me)
[![Go Report Card](https://goreportcard.com/badge/github.com/tomwright/dasel/v3)](https://goreportcard.com/report/github.com/tomwright/dasel/v3)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tomwright/dasel)](https://pkg.go.dev/github.com/tomwright/dasel/v3)
![Test](https://github.com/TomWright/dasel/workflows/Test/badge.svg)
![Build](https://github.com/TomWright/dasel/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/TomWright/dasel/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/dasel)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
![GitHub Downloads](https://img.shields.io/github/downloads/TomWright/dasel/total)
![Homebrew Formula Downloads](https://img.shields.io/homebrew/installs/dy/dasel?label=brew%20installs)
![GitHub License](https://img.shields.io/github/license/TomWright/dasel)
[![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/TomWright/dasel?label=latest%20release)](https://github.com/TomWright/dasel/releases/latest)
[![Homebrew tag (latest by date)](https://img.shields.io/homebrew/v/dasel)](https://formulae.brew.sh/formula/dasel)

<div align="center">
    <img src="./daselgopher.png" alt="Dasel mascot" width="250"/>
</div>

# Dasel

Dasel（**Data-Select** 的简称）是一个命令行工具和库，用于查询、修改和转换 JSON、YAML、TOML、XML、CSV、KDL 等数据结构。

它提供了统一而强大的语法来遍历和更新数据，因此对开发者、DevOps 工程师以及数据加工类任务都非常有用。

---

## 功能特性

* **多格式支持**：JSON、YAML、TOML、XML、CSV、HCL、INI、KDL。
* **统一的查询语法**：用同一套选择器访问任意格式的数据。
* **查询与搜索**：以直观的语法提取值、列表或结构。
* **原地修改**：直接在结构化文件中更新、插入或删除值。
* **格式互转**：无缝地将 JSON → YAML、TOML → JSON 等格式相互转换。
* **易于脚本化**：通过简洁的 CLI 即可集成到 shell 脚本与管道中。
* **库形式可用**：可在 Go 项目中导入并使用。

---

## 安装

### Homebrew（macOS/Linux）

```sh
brew install dasel
```

### Go Install

```sh
go install github.com/tomwright/dasel/v3/cmd/dasel@master
```

### 预编译二进制

Linux、macOS 和 Windows 的预编译二进制可在 [Releases](https://github.com/TomWright/dasel/releases) 页面获取。

### 以上都不适用？

更多安装方式请见[安装文档](https://daseldocs.tomwright.me/getting-started/installation)。

---

## Shell 自动补全

Dasel 支持 Bash、Zsh、Fish 和 PowerShell 的 Tab 补全。为你所使用的 shell 生成并加载补全脚本：

```sh
# Bash
source <(dasel completion bash)

# Zsh
source <(dasel completion zsh)

# Fish
dasel completion fish | source

# PowerShell
dasel completion powershell | Out-String | Invoke-Expression
```

关于永久安装补全脚本的方法，请见[Shell 补全文档](https://daseldocs.tomwright.me/getting-started/shell-completion)。

## Man 手册页

生成并查看 man 手册页：

```sh
dasel man | man -l -
```

---

## 基本用法

### 选择值

默认情况下，Dasel 会对最后一个选择器求值并打印结果。

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'foo.bar'
# Output: "baz"
```

### 修改值

内联更新值：

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'foo.bar = "bong"'
# Output: "bong"
```

使用 `--root` 在修改后输出完整文档：

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json --root 'foo.bar = "bong"'
# Output:
{
  "foo": {
    "bar": "bong"
  }
}
```

根据当前值更新目标值：

```sh
echo '[1,2,3,4,5]' | dasel -i json --root 'each($this = $this*2)'
# Output:
[
    2,
    4,
    6,
    8,
    10
]
```

### 格式转换

```sh
cat data.json | dasel -i json -o yaml
```

### 递归下降（`..`）

在所有嵌套的对象和数组中查找匹配的键或索引。

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json '..bar'
# Output:
[
    "baz"
]

```

### 搜索（`search`）

在结构的任意位置查找所有满足条件的值。

```sh
echo '{"foo": {"bar": "baz"}}' | dasel -i json 'search(bar == "baz")'
# Output:
[
    {
        "bar": "baz"
    }
]

```

---

## 文档

完整文档见 [daseldocs.tomwright.me](https://daseldocs.tomwright.me)。

---

## 贡献

欢迎贡献！详情请见 [CONTRIBUTING.md](./CONTRIBUTING.md)。

---

## 许可证

MIT 许可证。详见 [LICENSE](./LICENSE)。

## Star 增长趋势

[![Stargazers over time](https://starchart.cc/TomWright/dasel.svg)](https://starchart.cc/TomWright/dasel)
