package main

import (
	"bytes"
	"context"
	"kava/internal/configuration"
	"kava/internal/database"
	"kava/internal/database/compute"
	"kava/internal/database/storage"
	"kava/internal/database/storage/engine/in_memory"
	initialization "kava/internal/initalization"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(data)
	cfg, err := configuration.Load(reader)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := initialization.CreateLogger(cfg.Logging)
	if err != nil {
		log.Fatal(err)
	}

	compute, err := compute.NewCompute(logger)
	if err != nil {
		log.Fatal(err)
	}
	engine, err := in_memory.NewEngine(logger)
	if err != nil {
		log.Fatal(err)
	}

	wal, err := initialization.CreateWAL(cfg.WAL, logger)
	if err != nil {
		log.Fatal("failed to initialize wal")
	}

	wal.Start(ctx)
	
	storage, err := storage.NewStorage(engine, wal, logger)
	if err != nil {
		log.Fatal(err)

	}
	database, err := database.NewDatabase(compute, storage, logger)
	if err != nil {
		log.Fatal(err)
	}

	servers := initialization.NewServers(cfg, database, logger)

	var wg sync.WaitGroup
	for _, server := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			server.Start(ctx)
		}()
	}
	wg.Wait()
}
