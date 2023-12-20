package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"urlShortener/internal/config"
	"urlShortener/internal/gRPC/gRPCServer"
	"urlShortener/internal/http/httpServer"
	route "urlShortener/internal/http/htttpHandlers/router"
	"urlShortener/internal/lib/linkShortening"
	"urlShortener/internal/lib/linkShortening/hashByID"
	"urlShortener/internal/service"
	"urlShortener/internal/storage"
	"urlShortener/internal/storage/inMemmory"
	"urlShortener/internal/storage/postgres"
	_ "urlShortener/internal/storage/postgres"
	"urlShortener/pkg/logger"
)

const postgresStorage = "postgres"
const inMemoryStorage = defaultStorageType

func main() {
	flagsData := parseFlags()

	cfg, err := config.MustParseConfig(flagsData.cfgPath)
	if err != nil {
		log.Fatalf("cfg error: %v", err)
	}

	appLogger := logger.New()
	if err != nil {
		log.Fatalf("log error: %v", err)
	}

	appLogger.Infof("storage: %s", flagsData.storageType)

	ctx, final := context.WithCancel(context.Background())

	var db storage.Storager
	var hashGen linkShortening.Hasher

	switch flagsData.storageType {
	case postgresStorage:
		pq, err := postgres.New(ctx, &cfg.Postgres, appLogger)
		if err != nil {
			appLogger.Fatalf("can't init storage: %v", err)
		}
		maxID, err := pq.MaxID()
		if err != nil {
			appLogger.Fatalf("can't get maxID: %v", err)
		}
		hashGen = hashByID.New(maxID)
		db = pq
	case inMemoryStorage:
		db = inMemmory.New()
		hashGen = hashByID.New(0)
	default:
		appLogger.Fatalf("wrong storage type")
	}

	urlShortener := service.New(db, hashGen)

	router := route.New(appLogger, urlShortener)

	appLogger.Info("starting gRPCServer")

	srvGRPC := gRPCServer.New(appLogger)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = srvGRPC.Run(ctx, cfg.GRPCAddr, urlShortener)
		if err != nil {
			appLogger.Fatalf("can't run grpc %v: ", err)
		}
		wg.Done()
	}()

	srv := httpServer.New(ctx, cfg.HTTPServer, router, appLogger)
	appLogger.Debug(cfg.HTTPServer)
	appLogger.Info("starting HTTPServer")

	wg.Add(1)
	go func(){
		srv.Run()
		wg.Done()
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case <-signalCh:
		appLogger.Info("Received interrupt signal, shutting down")
	}

	final()
	wg.Wait()

	appLogger.Info("Server stopped gracefully")
}
