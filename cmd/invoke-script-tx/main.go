package main

import (
	"context"
	"net/http"
	"time"

	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func main() {
	// Create sender's public key from BASE58 string
	sender, err := crypto.NewPublicKeyFromBase58("<your-public-key>")
	if err != nil {
		panic(err)
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58("<your-private-key>")
	if err != nil {
		panic(err)
	}

	// Create script's address
	a, err := proto.NewAddressFromString("<script's address")
	if err != nil {
		panic(err)
	}

	// Create Function Call that will be passed to the script
	fc := proto.FunctionCall{
		Name: "foo",
		Arguments: proto.Arguments{
			proto.IntegerArgument{
				Value: 12345,
			},
			proto.BooleanArgument{
				Value: true,
			},
		},
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	// Fee asset is Waves
	waves := proto.OptionalAsset{Present: false}

	// New InvokeScript Transaction
	tx, err := proto.NewUnsignedInvokeScriptV1('T', sender, a, fc, proto.ScriptPayments{}, waves, 500000, uint64(ts))
	if err != nil {
		panic(err)
	}

	// Sing the transaction with the private key
	err = tx.Sign(sk)

	// Create new HTTP client to send the transaction to public TestNet nodes
	client, err := client.NewClient(client.Options{BaseUrl: "https://testnodes.wavesnodes.com", Client: &http.Client{}})
	if err != nil {
		panic(err)
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send the transaction to the network
	_, err = client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		panic(err)
	}
}
