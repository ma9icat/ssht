# SSH 任务执行工具 (ssht)

一个轻量级的在多个SSH主机上执行命令的工具。

简体中文 | [English](./README.md)

## 功能特性

- 在多个SSH主机上执行命令
- 使用`--nodes`参数过滤节点
- 可配置的日志格式(text/json)
- 使用`--debug`标志启用调试模式

## 安装

### 从源码安装
```bash
go build -o ssht main.go
```

### 使用Docker
```bash
docker build -t ssht .
```

## 使用说明

1. 在`config.toml`中配置您的主机
2. 运行命令:

```bash
# 基础用法
./ssht --command "hostname"

# 在指定节点上运行
./ssht --command "hostname" --nodes node1,node2

# 调试模式+JSON日志
./ssht --command "hostname" --nodes node1 --debug --log-format json

# Docker使用
docker run -v $(pwd)/config.toml:/app/config.toml ssht --command "hostname"
```

## 配置

编辑`config.toml`来配置您的SSH主机和分组。

## 日志

支持的日志格式:
- `text` (默认，彩色输出)
- `json` (结构化日志)

日志选项:
- `--log-format`: 指定格式(text/json)
- `--log-file`: 将日志写入指定文件(默认: 控制台)

示例:
```bash
# 将JSON日志写入文件
./ssht --command "hostname" --log-format json --log-file logs.json

# 将文本日志写入文件
./ssht --command "hostname" --log-file logs.txt
```