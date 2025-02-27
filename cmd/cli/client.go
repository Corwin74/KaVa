package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"kava/internal/database/client"
	"os"
	"syscall"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"go.uber.org/zap"
)

const (
	defaultBufferSize = 4 << 10
	defaultIdleTimeout = 60 * time.Second
)

func main() {
	address := flag.String("address", "localhost:8080", "Address of the KaVa")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	if *idleTimeout == 0 {
		*idleTimeout = defaultIdleTimeout
	}

	logger, _ := zap.NewProduction()
	_, err := bytefmt.ToBytes(*maxMessageSizeStr)

	if err != nil {
		logger.Fatal("failed to parse max message size", zap.Error(err))
	}

	reader := bufio.NewReader(os.Stdin)
	client, err := client.NewTCPClient(*address, defaultBufferSize, *idleTimeout)

	if err != nil {
		logger.Fatal("failed to create client", zap.Error(err))
	}

	for {
		fmt.Print("[kava] > ")
		request, err := reader.ReadString('\n')
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := client.Send([]byte(request))
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}
