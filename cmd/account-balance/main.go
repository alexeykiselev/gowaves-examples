package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tyler-smith/go-bip39"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func main() {
	// This is a simple example of how to create a new key pair. Get an address for it. And request a balance.

	// Generate a new random entropy bytes
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		panic(err)
	}

	// Make a mnemonic seed phrase out of the entropy
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		panic(err)
	}
	fmt.Println("Seed:", mnemonic)

	// Generate a key pair for the seed phrase
	// The secret key is not used later, so omit it (first return value)
	_, pk, err := crypto.GenerateKeyPair([]byte(mnemonic))
	if err != nil {
		panic(err)
	}

	// Make an address for the public key and Testnet
	addr, err := proto.NewAddressFromPublicKey(proto.TestNetScheme, pk)
	if err != nil {
		panic(err)
	}
	fmt.Println("Address:", addr.String())

	// Initialize an API client with public node's URL on Testnet
	cl, err := client.NewClient(client.Options{BaseUrl: "https://nodes-testnet.wavesnodes.com", Client: &http.Client{}})
	if err != nil {
		panic(err)
	}

	// Context to cancel the request execution on timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	balance, _, err := cl.Addresses.Balance(ctx, addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Balance (wavelets):", balance.Balance)
}
