package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/Shopify/sarama"
	"github.com/divpro/transactions-example/internal/config"
	"github.com/divpro/transactions-example/internal/consumer"
	"github.com/divpro/transactions-example/internal/consumer/handlers"
	"github.com/divpro/transactions-example/internal/repository"
	"github.com/divpro/transactions-example/internal/service"
	"github.com/divpro/transactions-example/pkg/entity"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "Configuration file name")
	flag.Parse()
}

func main() {
	textHandler := slog.NewTextHandler(os.Stdout)
	logger := slog.New(textHandler)

	f, err := os.Open(configPath)
	if err != nil {
		logger.Error("open configuration file", err, configPath)
		return
	}
	var conf config.Config
	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		logger.Error("parse configuration file", err, configPath)
		return
	}

	db, err := sql.Open("pgx", conf.DB.DSN())
	if err != nil {
		logger.Error("open db", err, conf.DB.DSN())
		return
	}
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	depositsRepo := repository.NewDeposits(db)
	transactionsRepo := repository.NewTransactions(db)
	usersRepo := repository.NewUsers(db)
	fin := service.NewFinance(depositsRepo, transactionsRepo, usersRepo)

	saramaConf := sarama.NewConfig()
	saramaConf.ClientID = "transactions-example"
	saramaConf.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	saramaConf.Consumer.Offsets.Initial = sarama.OffsetOldest

	depositGroup, err := sarama.NewConsumerGroup(conf.Kafka.Brokers, "deposits", saramaConf)
	if err != nil {
		logger.Error("create deposits consumer group", err, conf.Kafka.Brokers)
		return
	}
	depositHandler := handlers.NewDeposit(fin)
	depositConsumer := consumer.NewConsumer[entity.Deposit](logger, depositHandler)

	transactionsGroup, err := sarama.NewConsumerGroup(conf.Kafka.Brokers, "transactions", saramaConf)
	if err != nil {
		logger.Error("create transactions consumer group", err, conf.Kafka.Brokers)
		return
	}
	transactionsHandler := handlers.NewTransaction(fin)
	transactionsConsumer := consumer.NewConsumer[entity.Transaction](logger, transactionsHandler)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			if err := depositGroup.Consume(ctx, []string{"deposits"}, depositConsumer); err != nil {
				logger.Error("consume deposits", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			if err := transactionsGroup.Consume(ctx, []string{"transactions"}, transactionsConsumer); err != nil {
				logger.Error("consume transactions", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm
	logger.Info("terminating")
	cancel()
	wg.Wait()
	if err = depositGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	if err = transactionsGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
