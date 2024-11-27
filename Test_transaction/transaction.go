package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	ganacheURL = "http://localhost:7545"
)


// note :- update senders and theri private key to ganache senders and their private keys

var (
	// Senders
	sender1 = common.HexToAddress("0xdD93e92dc32d0B2F51430b0e6dA29BDd01AF68D6")
	sender2 = common.HexToAddress("0xC22c7f8bA7dE381A299ee4EB3a11E1316525ce45")
	sender3 = common.HexToAddress("0x91A9CeF0099DF1D2eA4F4B825ac06B506dfDbe07")

	// Recipients (accounts 4, 5, 6)
	recipients = []common.Address{
		common.HexToAddress("0xe53e48483b29A88c97Cb7B52120D0d59A86AF8E3"), // Account 4
		common.HexToAddress("0x67775427510E59072bA5752763Ba14B57F22dA46"), // Account 5
		common.HexToAddress("0x2d26457305a3Bf47dEb96E79f42F8a5654bD470b"), // Account 6
	}

	// Private keys for each sender
	privateKeys = []string{
		"99102aee559622da2dbac6cc042e797cdf433c4b70c983261d26452c7ddd15f6", // Sender 1
		"09d2e75941e3862710388f6efe118866a61c9d0f7fe41340a64edc19247b74d3", // Sender 2
		"19d8e06047c38f7cf2327b8dbcc29cbba850f39e62fa99a187691c95c76f6a09", // Sender 3
	}
)

func getBalance(client *ethclient.Client, address common.Address) (*big.Float, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}

	balanceEth := new(big.Float)
	balanceEth.SetString(balance.String())
	return new(big.Float).Quo(balanceEth, big.NewFloat(1e18)), nil
}

func processSenderTransactions(client *ethclient.Client, senderAddress common.Address, privateKeyHex string, recipient common.Address) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	senderBalance, err := getBalance(client, senderAddress)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %v", err)
	}
	fmt.Printf("\nSender %s initial balance: %v ETH\n", senderAddress.Hex(), senderBalance)

	// 0.1 ETH in Wei
	amount := new(big.Int).Mul(big.NewInt(100000000000000000), big.NewInt(1))

	nonce, err := client.PendingNonceAt(context.Background(), senderAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	recipientBalance, err := getBalance(client, recipient)
	if err != nil {
		return fmt.Errorf("failed to get recipient balance: %v", err)
	}

	fmt.Printf("Sending 0.1 ETH from %s to %s\n", senderAddress.Hex(), recipient.Hex())
	fmt.Printf("Recipient balance before: %v ETH\n", recipientBalance)

	txHash, err := sendTransaction(client, privateKey, recipient, amount, nonce)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction successful! Hash: %s\n", txHash)

	time.Sleep(2 * time.Second)

	newBalance, err := getBalance(client, recipient)
	if err != nil {
		return fmt.Errorf("failed to get new balance: %v", err)
	}
	fmt.Printf("Recipient balance after: %v ETH\n", newBalance)

	return nil
}
func sendTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, amount *big.Int, nonce uint64) (string, error) {
	ctx := context.Background()

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to suggest gas price: %v", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Create transaction data
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasFeeCap: gasPrice,
		GasTipCap: big.NewInt(1000000000),
		Gas:       gasLimit,
		To:        &to,
		Value:     amount,
		Data:      nil,
	}

	// Create and sign transaction
	signedTx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(chainID), txData)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	return signedTx.Hash().Hex(), nil
}
func main() {
	client, err := ethclient.Dial(ganacheURL)
	if err != nil {
		log.Fatalf("Failed to connect to Ganache: %v", err)
	}
	defer client.Close()

	senders := []common.Address{sender1, sender2, sender3}

	// Process transactions for each sender to their corresponding recipient
	for i := 0; i < len(senders); i++ {
		err := processSenderTransactions(client, senders[i], privateKeys[i], recipients[i])
		if err != nil {
			log.Printf("Error processing transactions for sender %s: %v", senders[i].Hex(), err)
		}
		time.Sleep(1 * time.Second)
	}

	// Print final balances
	fmt.Println("\nFinal Balances:")
	for i, sender := range senders {
		balance, err := getBalance(client, sender)
		if err != nil {
			log.Printf("Failed to get final balance for sender %s: %v", sender.Hex(), err)
			continue
		}
		fmt.Printf("Sender %d final balance: %v ETH\n", i+1, balance)
	}
}
