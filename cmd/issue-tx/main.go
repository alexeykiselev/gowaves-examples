package main

import (
	"context"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const waves = 100_000_000

func main() {
	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58("<your-private-key>")
	if err != nil {
		panic(err)
	}

	// Generate public key from secret key
	pk := crypto.GeneratePublicKey(sk)

	// Current time in milliseconds
	ts := uint64(time.Now().UnixMilli())

	// New Issue Transaction
	tx := proto.NewUnsignedIssueWithProofs(3, pk, "TestAsset", "Test Asset for Testnet",
		1000_00, 2, true, nil, ts, 1*waves)

	// Sing the transaction with the private key
	err = tx.Sign(proto.TestNetScheme, sk)
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

	// Send the transaction to the network
	_, err = cl.Transactions.Broadcast(ctx, tx)
	if err != nil {
		panic(err)
	}
}
