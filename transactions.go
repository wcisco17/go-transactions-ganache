package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Handling all errors
func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func main() {
	// HTTP server running on ganache
	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	checkError(err)

	// Private Key for buyers acccount
	privateKey, err := crypto.HexToECDSA("368e13a23e4ed34e69eb2983418c249664d187ffc85bee8180077bb4930f8b27")
	checkError(err)

	// Buyer's account address
	buyer := common.HexToAddress("0x76E8e249CDE25fEE6e5D73765604f9A764CF62Ac")

	// Checking buyers balance (fromAccount)
	{
		fromAccount := common.HexToAddress("0x76E8e249CDE25fEE6e5D73765604f9A764CF62Ac")
		balance, err := client.BalanceAt(context.Background(), fromAccount, nil)
		checkError(err)
		fmt.Println("Buyer Balance: ", balance)
	}

	// Calling pending nounce - Preventing double spending attack
	nonce, err := client.PendingNonceAt(context.Background(), buyer)
	checkError(err)
	value := new(big.Int)

	// In wei: 100000000000000000, convert ether: 0.01
	value.SetString("100000000000000000", 10)
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	checkError(err)

	// Painter's account address
	painter := common.HexToAddress("0x7922547d2E294fe16e329b3336491410b05B236d")

	// Checking The Painter's balance (toAccount)
	{
		toAccount := common.HexToAddress("0x7922547d2E294fe16e329b3336491410b05B236d")
		balance, err := client.BalanceAt(context.Background(), toAccount, nil)
		checkError(err)
		fmt.Println("Painter Balance: ", balance)
	}

	var data []byte

	// Executing Transaction
	tx := types.NewTransaction(nonce, painter, value, gasLimit, gasPrice, data)

	// Signing Transactions if err print it out
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	checkError(err)

	// Sending transaction amount in case there is not enough checkError
	err = client.SendTransaction(context.Background(), signedTx)
	checkError(err)

	// Print out the transaction hash if success
	fmt.Printf("Transaction Success! sent: %s", signedTx.Hash().Hex())
}
