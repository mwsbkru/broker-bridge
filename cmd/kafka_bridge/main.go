package main

import (
	"context"
	"github.com/mwsbkru/broker-bridge/internal/common"
	cfg "github.com/mwsbkru/broker-bridge/internal/config"
	"github.com/mwsbkru/broker-bridge/internal/kafka"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
)

const CONFIG_FILE_NAME = "config.yaml"

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)

	go func() {
		select {
		case <-ctx.Done():
			log.Println("Waiting bridge to terminate...")
		}
	}()

	var bridge common.Bridge
	bridgeConfig := loadConfig()

	bridge = kafka.NewBridge(bridgeConfig.Bridge.Kafka)
	bridge.Run(ctx)

	cancelFunc()
	log.Println("Bridge terminated!")
}

func loadConfig() cfg.Config {
	var bridgeConfig cfg.Config
	yamlFile, err := os.ReadFile(CONFIG_FILE_NAME)
	if err != nil {
		log.Fatalf("Load config error: %v.", err)
	}

	err = yaml.Unmarshal(yamlFile, &bridgeConfig)
	if err != nil {
		log.Fatalf("Unmarshal error: %v.", err)
	}

	return bridgeConfig
}
