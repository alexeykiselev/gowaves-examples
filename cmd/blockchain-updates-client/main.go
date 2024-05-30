package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	bu "github.com/wavesplatform/gowaves/pkg/grpc/generated/waves/events/grpc"
)

func main() {
	// Create a context that is canceled when the user interrupts the program.
	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer done()

	// Connect to blockchain updates node.
	conn, err := grpc.NewClient("nodes-testnet.wavesnodes.com:6881",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(10<<20))) // Enable messages up to 10MB
	if err != nil {
		panic(fmt.Errorf("failed to connect to node: %w", err))
	}

	// Create a new subscriber.
	c := bu.NewBlockchainUpdatesApiClient(conn)
	req := &bu.SubscribeRequest{
		FromHeight: int32(3_100_000), // We want to receive updates starting from block 3,100,000.
	}
	stream, err := c.Subscribe(ctx, req)
	if err != nil {
		panic(fmt.Errorf("failed to subscribe to blockchain updates: %w", err))
	}
	var event = new(bu.SubscribeEvent)
	for err = stream.RecvMsg(event); err == nil; err = stream.RecvMsg(event) {
		// Get event height.
		h := event.GetUpdate().GetHeight()
		// Switching by event type. Is it a block append or a rollback?
		switch u := event.GetUpdate().Update.(type) {
		case *events.BlockchainUpdated_Append_: // Extract block generator form block append event.
			if b := u.Append.GetBlock(); b != nil {
				generatorPublicKey, err := crypto.NewPublicKeyFromBytes(b.GetBlock().GetHeader().GetGenerator())
				if err != nil {
					panic(fmt.Errorf("failed to create public key from bytes: %w", err))
				}
				generatorAddress, err := proto.NewAddressFromPublicKey(proto.MainNetScheme, generatorPublicKey)
				if err != nil {
					panic(fmt.Errorf("failed to create address from public key: %w", err))
				}
				fmt.Printf("%d: Generator: %s\n", h, generatorAddress.String())
			}
		case *events.BlockchainUpdated_Rollback_: // Ignore rollbacks.
		default:
			panic(fmt.Errorf("unsupported event type %T at height %d", event.GetUpdate().Update, h))
		}
	}

	<-ctx.Done()
}
