package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"stockbit/config"
	"stockbit/domain"
	httpHandler "stockbit/handler/http"
	processorHandler "stockbit/handler/processor"
)

func main() {
	command := flag.String("command", "", "service command: http/processor")
	processor := flag.String("processor", "", "processor command: balance/above_threshold")
	flag.Parse()

	configJSONFile, err := os.Open("config/config.json")
	if err != nil {
		panic(err)
	}
	defer configJSONFile.Close()

	configBytes, err := ioutil.ReadAll(configJSONFile)
	if err != nil {
		panic(err)
	}

	var cfg config.Config
	err = json.Unmarshal(configBytes, &cfg)
	if err != nil {
		panic(err)
	}

	tmc := goka.NewTopicManagerConfig()
	tmcCfg := cfg.Goka.TopicManagerConfig
	tmc.Table.Replication = tmcCfg.TableReplication
	tmc.Stream.Replication = tmcCfg.StreamReplication

	brokers := cfg.Goka.Brokers

	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	for _, topic := range cfg.Goka.Topics {
		err = tm.EnsureStreamExists(topic.Name, topic.Partition)
		if err != nil {
			panic(err)
		}
	}

	if command == nil {
		panic("missing command")
	}

	switch *command {
	case "http":
		emitters := make(map[string]*goka.Emitter)
		for _, topic := range cfg.Goka.Topics {
			topicName := topic.Name

			emitters[topicName], err = goka.NewEmitter(brokers, goka.Stream(topicName), new(codec.String))
			if err != nil {
				panic(err)
			}
		}

		defer func() {
			for _, emitter := range emitters {
				_ = emitter.Finish()
			}
		}()

		viewers := make(map[string]*goka.View)

		viewers[domain.BalanceGroup], err = goka.NewView(
			brokers,
			goka.GroupTable(domain.BalanceGroup),
			new(domain.BalanceCodec),
		)
		if err != nil {
			panic(err)
		}

		viewers[domain.AboveThresholdGroup], err = goka.NewView(
			brokers,
			goka.GroupTable(domain.AboveThresholdGroup),
			new(domain.AboveThresholdCodec),
		)
		if err != nil {
			panic(err)
		}

		h := httpHandler.NewHTTPHandler(emitters, viewers)

		mux := http.NewServeMux()
		mux.HandleFunc("/deposit", h.Deposit)

		server := new(http.Server)
		server.Addr = fmt.Sprintf("%s:%s", cfg.Server.Address, cfg.Server.Port)
		server.Handler = mux

		for viewName, view := range viewers {
			go func(viewName string, view *goka.View) {
				errView := view.Run(context.Background())
				if errView != nil {
					log.Printf("error running %s: %s", viewName, errView)
					return
				}
			}(viewName, view)
		}

		if err = server.ListenAndServe(); err != nil {
			panic(err)
		}

	case "processor":
		if processor == nil {
			panic("missing processor command")
		}

		var ph domain.Processor

		switch *processor {
		case "balance":
			balanceProcessor := processorHandler.NewBalanceProcessor(cfg)
			ph = processorHandler.NewProcessorHandler(
				domain.DepositsTopic,
				domain.BalanceGroup,
				tmc,
				brokers,
				new(domain.BalanceCodec),
				balanceProcessor.Handle,
			)
		case "above_threshold":
			aboveThresholdProcessor := processorHandler.NewAboveThresholdProcessor(cfg)
			ph = processorHandler.NewProcessorHandler(
				domain.DepositsTopic,
				domain.AboveThresholdGroup,
				tmc,
				brokers,
				new(domain.AboveThresholdCodec),
				aboveThresholdProcessor.Handle,
			)
		default:
			panic(fmt.Errorf("unsupported processor command: %s", *processor))
		}

		ph.Run(context.Background())
	default:
		panic(fmt.Errorf("unsupported command: %s", *command))
	}
}
