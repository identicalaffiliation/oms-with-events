package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/broker"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/config"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/database"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/infrastructure/logger"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/repository"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/transport/rest"
	"github.com/identicalaffiliation/oms-with-events/order-service/internal/usecase"
	"github.com/identicalaffiliation/oms-with-events/order-service/pkg/dispatcher"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "c", "config.yml", "path to config file")
	flag.Parse()

	cfg, err := config.NewOrderServiceConfig(configPath)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to init config", err)
		os.Exit(1)
	}

	slogger, err := logger.NewSLogger(cfg)
	if err != nil {
		fmt.Fprintln(os.Stdout, "failed to init logger", err)
		os.Exit(1)
	}

	pool, err := database.NewPool(cfg)
	if err != nil {
		slogger.Error("failed to init new postgres pool", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err := pool.Close(); err != nil {
			slogger.Error("failed to close postgres pool", "error", err)
		}
	}()

	slogger.Debug("infrastructure added")

	ordersRepository := repository.NewOrdersRepository(pool)
	eventsRepository := repository.NewEventsRepository(pool)
	ordersUsecase := usecase.NewOrdersUsecase(ordersRepository, eventsRepository, pool, slogger)

	producer := broker.NewProducer(cfg.KafkaConfig.Brokers, cfg.KafkaConfig.Topic)
	dispatcher := dispatcher.NewDispatcher(producer, eventsRepository, cfg.DispatcherConfig.WorkersCount,
		cfg.DispatcherConfig.BatchSize, cfg.DispatcherConfig.RetryCount,
		slogger, cfg.DispatcherConfig.ChillDuration, pool)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go dispatcher.Run(ctx)

	api := rest.NewOrderServiceAPI(ordersUsecase)

	server, shutdown := setupMux(api, cfg, slogger)

	signalsChan := make(chan os.Signal, 1)
	signal.Notify(signalsChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		slogger.Debug("starting server..")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slogger.Error("failed to listen port", "error", err)
		}
	}()

	<-signalsChan

	shutdownContext, cancel := context.WithTimeout(context.Background(), cfg.ServerConfig.ShutdownTimeout)
	defer cancel()

	shutdown(shutdownContext)

	slogger.Debug("server stopped gracefully..")
}

func setupMux(
	api *rest.OrderServiceAPI,
	cfg *config.OMSGOrderServiceConfig,
	logger logger.Logger,
) (*http.Server, func(ctx context.Context)) {
	mux := chi.NewRouter()
	mux.Post("/orders", api.CreateOrder)
	mux.Get("/orders/my/{id}", api.GetMyOrders)

	serverAddress := fmt.Sprintf("%s:%d", cfg.ServerConfig.Host, cfg.ServerConfig.Port)

	server := &http.Server{
		Addr:         serverAddress,
		Handler:      mux,
		ReadTimeout:  cfg.ServerConfig.ReadTimeout,
		WriteTimeout: cfg.ServerConfig.WriteTimeout,
		IdleTimeout:  cfg.ServerConfig.IdleTimeout,
	}

	shutdownFunc := func(ctx context.Context) {
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown server", "error", err)
		}
	}

	return server, shutdownFunc
}
