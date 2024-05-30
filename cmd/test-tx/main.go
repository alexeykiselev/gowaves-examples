package main

import (
	"log"
	"time"

	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

func main() {
	sk, pk, err := crypto.GenerateKeyPair([]byte("seed-seed-seed"))
	if err != nil {
		log.Fatal(err)
	}
	addr, err := proto.NewAddressFromPublicKey('T', pk)
	fc := proto.NewFunctionCall("unlockUserLock", proto.Arguments{
		&proto.StringArgument{
			Value: "pool",
		},
		&proto.StringArgument{
			Value: "user",
		}})
	ts := time.Now().Unix()
	tx := proto.NewUnsignedInvokeScriptWithProofs(
		byte(2),
		pk,
		proto.NewRecipientFromAddress(addr),
		fc,
		proto.ScriptPayments{},
		proto.OptionalAsset{},
		500000,
		uint64(ts),
	)
	err = tx.Sign('T', sk)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OK")
}
