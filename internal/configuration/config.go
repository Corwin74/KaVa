package configuration

import (
	"fmt"
	"io"
	"time"

	"code.cloudfoundry.org/bytefmt"

	"gopkg.in/yaml.v3"
)

const (
	defaultMaxConnections = 100
	defaultMaxMessageSize = 4 * bytefmt.KILOBYTE
	defaultIdleTimeout    = 60 * time.Second
	defaultHost           = "0.0.0.0"
	defaultPort           = 8080

	// Default logging values
	defaultLogLevel  = "info"
	defaultLogOutput = "stdout"

	// Default engine type
	defaultEngineType = "in_memory"
)

// Supported values constants
var (
	supportedLogLevels   = []string{"debug", "info", "warn", "error", "fatal"}
	supportedEngineTypes = []string{"in_memory", "persistent"}
)

// MessageSize - custom тип для срабатывания UnMarshal
type MessageSize int

// ServerConfig - базовый интерфейс для всех конфигураций серверов
type ServerConfig interface {
	// getType возвращает тип сервера
	getType() string

	// Validate проверяет корректность конфигурации
	// Validate() error

	// getName возвращает имя сервера
	getName() string
}

// ServerConfigs - slice для хранения конфигураций серверов
type ServerConfigs []ServerConfig

// Config -- корневая структура
type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Servers ServerConfigs  `yaml:"servers"`
	Logging *LoggingConfig `yaml:"logging"`
}

// EngineConfig -- раздел движка
type EngineConfig struct {
	Type string `yaml:"type"`
}

// LoggingConfig -- раздел логгера
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// BaseServer - базовая структура с общими полями
type BaseServer struct {
	Type string `yaml:"type"`
	Name string `yaml:"name"`
}

// TCPServerConfig - конфигурация TCP сервера
type TCPServerConfig struct {
	BaseServer     `yaml:",inline"`
	Port           int           `yaml:"port"`
	Host           string        `yaml:"host"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize MessageSize   `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

// Load -- загружает информацию из файла
func Load(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration: %w", err)
	}

	return &config, nil
}
