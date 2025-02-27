package configuration

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testCfgData = `engine:
  type: "in_memory"

servers:
  - type: tcp
    name: main
    port: 8087
    host: localhost
    max_connections: 1
    max_message_size: "2KB"
    idle_timeout: 5m

  - type: console
    name:  console-service

  - type: tcp
    name: hello-world

logging:
  level: "info"
  output: "output.log"
`

func TestLoad(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfgData string

		expectedCfg Config
	}{
		"load empty config": {cfgData: ""},
		"load config": {cfgData: testCfgData,
			expectedCfg: Config{
				Engine: &EngineConfig{
					Type: "in_memory",
				},
				Servers: []ServerConfig{
					&TCPServerConfig{
						BaseServer: BaseServer{
							Type: "tcp",
							Name: "main",
						},
						Port:           8087,
						Host:           "localhost",
						MaxConnections: 1,
						MaxMessageSize: 2048,
						IdleTimeout:    5 * time.Minute,
					},
					&ConsoleConfig{
						BaseServer: BaseServer{
							Type: "console",
							Name: "console-service",
						},
					},
					&TCPServerConfig{
						BaseServer: BaseServer{
							Type: "tcp",
							Name: "hello-world",
						},
						Port:           8080,
						Host:           "0.0.0.0",
						MaxConnections: 100,
						MaxMessageSize: 4096,
						IdleTimeout:    1 * time.Minute,
					},
				},
				Logging: &LoggingConfig{Level: "info", Output: "output.log"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader(test.cfgData)
			cfg, err := Load(reader)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedCfg, *cfg)
		})
	}
}