# Go Ethereum Sepolia Example

This repository contains a simple Go program for interacting with the Sepolia testnet using `go-ethereum`.

## 环境准备

1. 安装 Go
   - macOS: `brew install go`
   - 确保 `go version` 输出 Go 1.20+。

2. 安装 `go-ethereum` 模块
   - 程序依赖会在 `go mod tidy` 时自动下载。

3. 注册 Infura 并获取 Sepolia API Key
   - 登录 Infura，创建项目，获取 `INFURA_API_KEY`。
   - Sepolia RPC URL: `https://sepolia.infura.io/v3/<API_KEY>`。

4. 准备 Sepolia 测试账户
   - 在 MetaMask 或其他钱包中创建账户，并获取私钥
   - 请勿将私钥提交到版本控制。

## 运行示例

### 查询区块

```bash
export INFURA_API_KEY="your_infura_api_key"
go run main.go -mode query -block 2000000
```

### 发送交易

```bash
export INFURA_API_KEY="your_infura_api_key"
go run main.go -mode send -private-key "your_private_key" -to "0xRecipientAddress" -amount 0.01
```

如果使用自定义 RPC，请传入 `-endpoint`：

```bash
go run main.go -mode query -endpoint "https://sepolia.infura.io/v3/your_infura_api_key" -block 2000000
```
