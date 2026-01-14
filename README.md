# auth4x

## 快速开始

在生成 go-zero 风格（`style=go_zero`）的模板之前，请先完成环境安装与环境变量配置：

```bash
# 安装 Go（示例：macOS 可用 brew，其他系统请按官方文档安装）
# brew install go

# 建议设置 Go 环境变量（根据实际安装路径调整）
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

# 使环境变量生效（示例：bash）
# source ~/.bashrc

# 验证 Go 安装
go version
```

完成 Go 环境后，安装 `goctl`：

```bash
# 安装 goctl（建议放入 GOPATH/bin）
GO111MODULE=on go install github.com/zeromicro/go-zero/tools/goctl@latest

# 验证安装
which goctl
goctl --version
```

接下来即可使用 `goctl` 按 go-zero 风格生成模板：

```bash
# 仅示例：根据需要替换实际的 api/项目参数
# goctl api new auth4x --style=go_zero
```

如需生成基础的 logto 版本，可在模板生成后继续补充相关 API 与服务逻辑。
