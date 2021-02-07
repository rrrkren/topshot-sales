package main

import (
	"context"
	"fmt"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/rrrkren/topshot-sales/topshot"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Config struct {
	MomentEventDepth int    `mapstructure:"MOMENT_EVENT_DEPTH"`
	MomentEventType  string `mapstructure:"MOMENT_EVENT_TYPE"`
	OnflowHostname   string `mapstructure:"ONFLOW_HOSTNAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("default")
	viper.SetConfigType("env")
	viper.BindEnv("MOMENT_EVENT_DEPTH")
	viper.BindEnv("MOMENT_EVENT_TYPE")
	viper.BindEnv("ONFLOW_HOSTNAME")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

func main() {

	config, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load config:", err)
	}
	fmt.Println(config)

	//To get from the toml file or env var

	// connect to flow
	flowClient, err := client.New(config.OnflowHostname, grpc.WithInsecure())
	handleErr(err)
	err = flowClient.Ping(context.Background())
	handleErr(err)

	// fetch latest block
	latestBlock, err := flowClient.GetLatestBlock(context.Background(), false)
	handleErr(err)
	fmt.Println("current height: ", latestBlock.Height)

	// fetch block events of defined type and depth, defaults are topshot Market.MomentPurchased events for the past 500 blocks
	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
		Type:        config.MomentEventType,
		StartHeight: latestBlock.Height - uint64(config.MomentEventDepth),
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
