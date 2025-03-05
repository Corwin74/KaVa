package initialization

import (
	"context"
	"kava/internal/configuration"
	"kava/internal/database"
	"kava/internal/database/server"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
)

// Server -- интерфейс сервера
type Server interface {
	Start(context.Context)
}

func NewServers(cfg *configuration.Config, database *database.Database, logger *zap.Logger) []Server {
	var wg sync.WaitGroup

	servers := make([]Server, 0, len(cfg.Servers))

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
				servers = append(servers, tcpServer)
			case *configuration.ConsoleConfig:
				console, err := server.NewConsole(os.Stdin, os.Stdout, database, logger)
				if err != nil {
					log.Printf("failed to create console: %v", err)
					return
				}
				servers = append(servers, console)
			}
		}()
	}
	wg.Wait()

	return servers

}
