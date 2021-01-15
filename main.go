package main

import (
	"context"
	"fmt"
	"github.com/rrrkren/topshot-sales/topshot"

	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// connect to flow
	flowClient, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	handleErr(err)
	err = flowClient.Ping(context.Background())
	handleErr(err)

	// fetch latest block
	latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	handleErr(err)
	fmt.Println("current height: ", latestBlock.Height)

	// fetch block events of topshot Market.MomentPurchased events for the past 1000 blocks
	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        "A.c1e4f4f4c4257510.Market.MomentPurchased",
		StartHeight: latestBlock.Height - 500,
		EndHeight:   latestBlock.Height,
	})
	handleErr(err)

	for _, blockEvent := range blockEvents {
		for _, purchaseEvent := range blockEvent.Events {
			// loop through the Market.MomentPurchased events in this blockEvent
			e := topshot.MomentPurchasedEvent(purchaseEvent.Value)
			fmt.Println(e)
			saleMoment, err := topshot.GetSaleMomentFromOwnerAtBlock(flowClient, blockEvent.Height-1, *e.Seller(), e.Id())
			handleErr(err)
			fmt.Println(saleMoment)
			fmt.Printf("transactionID: %s, block height: %d\n",
				purchaseEvent.TransactionID.String(), blockEvent.Height)
			fmt.Println()
		}
	}
}
