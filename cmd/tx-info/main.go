package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

func main() {
	// Create transaction ID from Base58 string
	txID, err := crypto.NewDigestFromBase58("<base58-encoded-transaction-id")
	if err != nil {
		panic(err)
	}

	// Create new HTTP client to send the transaction to public TestNet nodes
	cl, err := client.NewClient(client.Options{BaseUrl: "https://nodes-testnet.wavesnodes.com", Client: &http.Client{}})
	if err != nil {
		panic(err)
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Request transaction info
	info, _, err := cl.Transactions.Info(ctx, txID)
	if err != nil {
		panic(err)
	}

	// Marshal transaction info to string and print it
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
