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

	storage, err := storage.NewStorage(engine, logger)
	if err != nil {
		log.Fatal(err)

	}
	database, err := database.NewDatabase(compute, storage, logger)
	if err != nil {
		log.Fatal(err)
	}

	servers := initialization.NewServers(cfg, database, logger)

	for _, server := range servers {
		server.Start(ctx)
	}

}
