# Go Ethereum Sepolia Example

This repository contains a simple Go program for interacting with the Sepolia testnet using `go-ethereum`, including smart contract deployment and interaction via `abigen`.

## 环境准备

1. 安装 Go
   - macOS: `brew install go`
   - 确保 `go version` 输出 Go 1.20+。

2. 安装 `go-ethereum` 模块
   - 程序依赖会在 `go mod tidy` 时自动下载。

3. 安装 Solidity 编译器
   - macOS: `brew install solidity`

4. 安装 abigen 工具
   ```bash
   go install github.com/ethereum/go-ethereum/cmd/abigen@v1.12.0
   ```

5. 注册 Infura 并获取 Sepolia API Key
   - 登录 Infura，创建项目，获取 `INFURA_API_KEY`。
   - Sepolia RPC URL: `https://sepolia.infura.io/v3/<API_KEY>`。

6. 准备 Sepolia 测试账户
   - 在 MetaMask 或其他钱包中创建账户，并获取私钥
   - 请勿将私钥提交到版本控制。

## 配置环境变量

创建 `.env` 文件（参考 `.env.example`）：

```bash
INFURA_API_KEY=your_infura_api_key
PRIVATE_KEY=your_private_key_hex
```

## 运行示例

### 查询区块

```bash
go run . -mode query -block 2000000
```

### 发送 ETH 交易

```bash
go run . -mode send -to "0xRecipientAddress" -amount 0.01
```

> `-private-key` 可选，不传时自动从 `.env` 的 `PRIVATE_KEY` 环境变量读取。

### 智能合约交互 (Counter)

#### 编译合约 & 生成 Go 绑定

```bash
# 编译 Solidity → ABI + 字节码
solc --abi --bin contracts/counter.sol -o contracts/

# 使用 abigen 生成 Go 绑定代码
abigen --abi=contracts/Counter.abi --bin=contracts/Counter.bin --pkg=contracts --out=contracts/counter.go
```

#### 部署合约

```bash
go run . -mode counter -action deploy
```

输出示例：
```
From: 0xAb49...
Balance: 9484264676690889815 wei
=== Deploying Counter Contract ===
Tx Hash: 0x...
Waiting for mining...
Contract deployed at: 0x...
Gas used: 151083
```

#### 查询计数

```bash
go run . -mode counter -action get -contract 0x<合约地址>
```

#### 递增计数

```bash
go run . -mode counter -action inc -contract 0x<合约地址>
```

#### 完整流程

```bash
# 1. 部署
go run . -mode counter -action deploy
# → 记下合约地址

# 2. 查初始值 (应为 0)
go run . -mode counter -action get -contract 0x<addr>

# 3. 递增
go run . -mode counter -action inc -contract 0x<addr>

# 4. 再查 (应为 1)
go run . -mode counter -action get -contract 0x<addr>
```

如果使用自定义 RPC，请传入 `-endpoint`：

```bash
go run . -mode query -endpoint "https://sepolia.infura.io/v3/your_infura_api_key" -block 2000000
```

## 项目结构

```
Go_Ethereum/
├── contracts/
│   ├── counter.sol          # Solidity 计数器合约
│   ├── Counter.abi          # 编译产物：ABI
│   ├── Counter.bin          # 编译产物：字节码
│   └── counter.go           # abigen 自动生成（不可手动编辑）
├── main.go                  # 程序入口 + query / send / counter 模式路由
├── counter_deploy.go        # 合约部署、读取、递增逻辑
├── .env.example             # 环境变量模板
├── .gitignore               # 忽略 .env 和二进制文件
├── go.mod
├── go.sum
└── README.md
```
