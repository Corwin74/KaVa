package main

import (
	"bytes"
	"context"
	"kava/internal/configuration"
	"kava/internal/database"
	"kava/internal/database/compute"
	"kava/internal/database/server"
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

	storage, err := storage.NewStorage(engine, logger)
	if err != nil {
		log.Fatal(err)

	}
	database, err := database.NewDatabase(compute, storage, logger)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, cfg := range cfg.Servers {
		wg.Add(1)
		go func() {
			defer wg.Done()

			switch cfg := cfg.(type) {
			case *configuration.TCPServerConfig:
				tcpServer, err := server.NewTCPServer(cfg, database, logger)
				if err != nil {
					log.Printf("failed to create tcp server: %v", err)
					return
				}
				tcpServer.Start(ctx)
			case *configuration.ConsoleConfig:
				console, err := server.NewConsole(os.Stdin, os.Stdout, database, logger)
				if err != nil {
					log.Printf("failed to create console: %v", err)
					return
				}
				console.Start(ctx)
			}
		}()
	}
	wg.Wait()
}
