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
	pk, err := crypto.NewPublicKeyFromBase58("<your-public-key>")
	if err != nil {
		panic(err)
	}

	// Create sender's private key from BASE58 string
	sk, err := crypto.NewSecretKeyFromBase58("<your-private-key>")
	if err != nil {
		panic(err)
	}

	// Current time in milliseconds
	ts := time.Now().Unix() * 1000

	// Create new alias with blockchain byte 'T' for TestNet
	alias := proto.NewAlias(proto.TestNetScheme, "testnetnode2")

	// New CreateAlias Transaction
	tx := proto.NewUnsignedCreateAliasWithSig(pk, *alias, 100_000, uint64(ts))

	// Sing the transaction with the private key
	err = tx.Sign(proto.TestNetScheme, sk)
	if err != nil {
		panic(err)
	}

	// Here the trickiest part, we have to convert the transaction to the request,
	// because the API accepts not the alias string representation, but alias value only
	req := client.AliasBroadcastReq{
		SenderPublicKey: pk,
		Fee:             tx.Fee,
		Timestamp:       tx.Timestamp,
		Signature:       *tx.Signature,
		Alias:           tx.Alias.Alias,
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
	_, _, err = cl.Alias.Broadcast(ctx, req)
	if err != nil {
		panic(err)
	}
}
