package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const sepoliaChainID = 11155111

func main() {
	loadEnv(".env")

	mode := flag.String("mode", "query", "operation mode: query or send")
	blockNumber := flag.String("block", "", "Sepolia block number to query (required for query mode)")
	txTo := flag.String("to", "", "recipient address for send mode")
	txAmount := flag.String("amount", "0.01", "amount in ETH for send mode")
	txPrivateKey := flag.String("private-key", "", "sender private key hex for send mode")
	rpcURL := flag.String("endpoint", "", "RPC endpoint; defaults to Sepolia Infura via INFURA_API_KEY env")
	flag.Parse()

	endpoint := getEndpoint(*rpcURL)
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatalf("failed to connect to Sepolia RPC: %v", err)
	}
	defer client.Close()

	switch strings.ToLower(*mode) {
	case "query":
		if *blockNumber == "" {
			log.Fatal("-block is required for query mode")
		}
		if err := queryBlock(client, *blockNumber); err != nil {
			log.Fatalf("query failed: %v", err)
		}
	case "send":
		if *txPrivateKey == "" || *txTo == "" {
			log.Fatal("-private-key and -to are required for send mode")
		}
		if err := sendTransaction(client, *txPrivateKey, *txTo, *txAmount); err != nil {
			log.Fatalf("send transaction failed: %v", err)
		}
	default:
		log.Fatalf("unknown mode: %s; valid values are query or send", *mode)
	}
}

func getEndpoint(explicit string) string {
	if explicit != "" {
		return explicit
	}
	apiKey := strings.TrimSpace(os.Getenv("INFURA_API_KEY"))
	if apiKey == "" {
		log.Fatal("INFURA_API_KEY env var not set and no -endpoint provided")
	}
	return fmt.Sprintf("https://sepolia.infura.io/v3/%s", apiKey)
}

func queryBlock(client *ethclient.Client, blockNum string) error {
	ctx := context.Background()
	bn, ok := new(big.Int).SetString(blockNum, 10)
	if !ok {
		return fmt.Errorf("invalid block number: %s", blockNum)
	}

	block, err := client.BlockByNumber(ctx, bn)
	if err != nil {
		return err
	}

	fmt.Println("=== Sepolia Block Info ===")
	fmt.Printf("Number: %d\n", block.NumberU64())
	fmt.Printf("Hash: %s\n", block.Hash().Hex())
	fmt.Printf("Parent Hash: %s\n", block.ParentHash().Hex())
	fmt.Printf("Timestamp: %d\n", block.Time())
	fmt.Printf("Transactions: %d\n", len(block.Transactions()))
	fmt.Printf("Gas Used: %d\n", block.GasUsed())
	fmt.Printf("Miner: %s\n", block.Coinbase().Hex())
	fmt.Printf("Difficulty: %s\n", block.Difficulty().String())
	return nil
}

func sendTransaction(client *ethclient.Client, privateKeyHex, toHex, amountEth string) error {
	ctx := context.Background()
	privateKey, err := loadPrivateKey(privateKeyHex)
	if err != nil {
		return err
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	value, err := parseEther(amountEth)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress(toHex)
	if toAddress == (common.Address{}) {
		return errors.New("invalid recipient address")
	}

	tx := types.NewTransaction(nonce, toAddress, value, 21000, gasPrice, nil)
	chainID := big.NewInt(sepoliaChainID)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	if err := client.SendTransaction(ctx, signedTx); err != nil {
		return err
	}

	fmt.Println("=== Transaction Sent ===")
	fmt.Printf("From: %s\n", fromAddress.Hex())
	fmt.Printf("To: %s\n", toAddress.Hex())
	fmt.Printf("Value: %s ETH\n", amountEth)
	fmt.Printf("Tx Hash: %s\n", signedTx.Hash().Hex())
	return nil
}

func loadPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	hexKey = strings.TrimPrefix(strings.TrimSpace(hexKey), "0x")
	if hexKey == "" {
		return nil, errors.New("private key is empty")
	}
	privBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key hex: %w", err)
	}
	return crypto.ToECDSA(privBytes)
}

func parseEther(amount string) (*big.Int, error) {
	f, ok := new(big.Float).SetString(strings.TrimSpace(amount))
	if !ok {
		return nil, fmt.Errorf("invalid ETH amount: %s", amount)
	}
	wei := new(big.Float).Mul(f, big.NewFloat(1e18))
	result := new(big.Int)
	wei.Int(result)
	return result, nil
}

func loadEnv(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return // no .env file is not an error
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key != "" && os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
}
