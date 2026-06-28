package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/madebyst/go-ethereum/contracts"
)

func deployCounter(client *ethclient.Client, privateKeyHex string) error {
	ctx := context.Background()

	privateKey, err := loadPrivateKey(privateKeyHex)
	if err != nil {
		return err
	}

	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}
	fmt.Printf("From: %s\n", fromAddress.Hex())
	fmt.Printf("Balance: %s wei\n", balance.String())

	chainID := big.NewInt(sepoliaChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}

	auth.GasLimit = 300000
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	address, tx, _, err := contracts.DeployContracts(auth, client)
	if err != nil {
		return fmt.Errorf("deploy contract: %w", err)
	}

	fmt.Println("=== Deploying Counter Contract ===")
	fmt.Printf("Tx Hash: %s\n", tx.Hash().Hex())
	fmt.Println("Waiting for mining...")

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("deployment reverted (status=0)")
	}

	fmt.Printf("Contract deployed at: %s\n", address.Hex())
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)
	return nil
}

func getCount(client *ethclient.Client, contractAddrHex string) error {
	address := common.HexToAddress(contractAddrHex)
	instance, err := contracts.NewContracts(address, client)
	if err != nil {
		return fmt.Errorf("bind contract: %w", err)
	}

	count, err := instance.GetCount(nil)
	if err != nil {
		return fmt.Errorf("call getCount: %w", err)
	}

	fmt.Println("=== Counter Contract ===")
	fmt.Printf("Contract: %s\n", address.Hex())
	fmt.Printf("Count: %s\n", count.String())
	return nil
}

func incrementCounter(client *ethclient.Client, contractAddrHex, privateKeyHex string) error {
	ctx := context.Background()

	privateKey, err := loadPrivateKey(privateKeyHex)
	if err != nil {
		return err
	}

	chainID := big.NewInt(sepoliaChainID)
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("create transactor: %w", err)
	}

	auth.GasLimit = 100000
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("get gas price: %w", err)
	}
	auth.GasPrice = gasPrice

	address := common.HexToAddress(contractAddrHex)
	instance, err := contracts.NewContracts(address, client)
	if err != nil {
		return fmt.Errorf("bind contract: %w", err)
	}

	tx, err := instance.Increment(auth)
	if err != nil {
		return fmt.Errorf("send increment tx: %w", err)
	}

	fmt.Println("=== Sending Increment ===")
	fmt.Printf("Contract: %s\n", address.Hex())
	fmt.Printf("Tx Hash: %s\n", tx.Hash().Hex())
	fmt.Println("Waiting for mining...")

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("increment reverted (status=0)")
	}

	fmt.Println("Increment confirmed!")
	fmt.Printf("Gas used: %d\n", receipt.GasUsed)

	count, err := instance.GetCount(nil)
	if err != nil {
		return fmt.Errorf("call getCount: %w", err)
	}
	fmt.Printf("New Count: %s\n", count.String())
	return nil
}
