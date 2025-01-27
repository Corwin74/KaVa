package main

import (
	"bufio"
	"context"
	"fmt"
	"kava/internal/database"
	"kava/internal/database/compute"
	"kava/internal/database/storage"
	"kava/internal/database/storage/engine/in_memory"
	"kava/internal/initalization"
	"log"
	"os"

	"go.uber.org/zap"
)


func main() {
	
	var res string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := initialization.CreateLogger()
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
	reader := bufio.NewReader(os.Stdin)
	for {
		query, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Failed to read query", zap.Error(err))
			continue
		}
		res = database.HandleQuery(ctx, query)

		fmt.Println(res)
	}
}
