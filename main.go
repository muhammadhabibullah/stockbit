package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/codec"
	"stockbit/domain"
	httpHandler "stockbit/handler/http"
	processorHandler "stockbit/handler/processor"
)

const (
	topic               = "deposits"
	balanceGroup        = "balance"
	aboveThresholdGroup = "aboveThresholdGroup"
)

func main() {
	command := flag.String("command", "", "service command: http/processor")
	processor := flag.String("processor", "", "processor command: balance/above_threshold")
	flag.Parse()

	tmc := goka.NewTopicManagerConfig()
	tmc.Table.Replication = 1
	tmc.Stream.Replication = 1

	brokers := []string{
		"localhost:9092",
	}
	tm, err := goka.NewTopicManager(brokers, goka.DefaultConfig(), tmc)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	err = tm.EnsureStreamExists("deposits", 3)
	if err != nil {
		panic(err)
	}

	emitter, err := goka.NewEmitter(brokers, topic, new(codec.String))
	if err != nil {
		panic(err)
	}
	defer emitter.Finish()

	if command == nil {
		panic("missing command")
	}

	switch *command {
	case "http":
		balanceView, err := goka.NewView(
			brokers,
			goka.GroupTable(balanceGroup),
			new(domain.BalanceCodec),
		)
		if err != nil {
			panic(err)
		}

		aboveThresholdView, err := goka.NewView(
			brokers,
			goka.GroupTable(aboveThresholdGroup),
			new(domain.AboveThresholdCodec),
		)
		if err != nil {
			panic(err)
		}

		h := httpHandler.NewHTTPHandler(emitter, balanceView, aboveThresholdView)

		mux := http.NewServeMux()
		mux.HandleFunc("/deposit", h.Deposit)

		server := new(http.Server)
		server.Addr = "localhost:8000"
		server.Handler = mux

		go func() {
			errView := balanceView.Run(context.Background())
			if errView != nil {
				log.Printf("error running balanceView: %s", errView)
				return
			}
		}()

		go func() {
			errView := aboveThresholdView.Run(context.Background())
			if errView != nil {
				log.Printf("error running balanceView: %s", errView)
				return
			}
		}()

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
			balanceProcessor := processorHandler.NewBalanceProcessor()
			ph = processorHandler.NewProcessorHandler(
				topic,
				balanceGroup,
				tmc,
				brokers,
				new(domain.BalanceCodec),
				balanceProcessor.Handle,
			)
		case "above_threshold":
			aboveThresholdProcessor := processorHandler.NewAboveThresholdProcessor()
			ph = processorHandler.NewProcessorHandler(
				topic,
				aboveThresholdGroup,
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
