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

// getType -- возращает тип сервера
func (b BaseServer) getType() string {
	return b.Type
}

// ConsoleConfig - конфигурация консольного сервера
type ConsoleConfig struct {
	BaseServer `yaml:",inline"`
}

func (c *ConsoleConfig) getType() string {
	return c.Type
}

func (c *ConsoleConfig) getName() string {
	return c.Name
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

func (t *TCPServerConfig) getType() string {
	return t.Type
}

func (t *TCPServerConfig) getName() string {
	return t.Name
}

// UnmarshalYAML реализует интерфейс yaml.Unmarshaler
func (m *MessageSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Сначала пробуем распарсить как строку
	var strSize string
	if err := unmarshal(&strSize); err == nil {
		// Если получилось распарсить как строку, конвертируем её
		size, err := bytefmt.ToBytes(strSize)
		if err != nil {
			return err
		}
		*m = MessageSize(size)
		return nil
	}

	// Если не получилось как строку, пробуем как число
	var size int
	if err := unmarshal(&size); err != nil {
		return err
	}

	*m = MessageSize(size)

	return nil
}

// UnmarshalYAML - кастомная десериализация для slice серверов
func (s *ServerConfigs) UnmarshalYAML(value *yaml.Node) error {
	// Проверяем, что получили sequence (массив) в YAML
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("expected sequence of servers, got %v", value.Kind)
	}

	// Временный slice для хранения результатов
	servers := make([]ServerConfig, 0, len(value.Content))

	// Обрабатываем каждый элемент в sequence
	for _, item := range value.Content {
		// Временная структура для определения типа сервера
		var baseServer struct {
			Type string `yaml:"type"`
		}

		if err := item.Decode(&baseServer); err != nil {
			return fmt.Errorf("failed to decode server type: %w", err)
		}

		// В зависимости от типа создаем соответствующую структуру
		var server ServerConfig
		switch baseServer.Type {
		case "console":
			var s ConsoleConfig
			if err := item.Decode(&s); err != nil {
				return fmt.Errorf("failed to decode console server: %w", err)
			}
			server = &s

		case "tcp":
			var s TCPServerConfig
			if err := item.Decode(&s); err != nil {
				return fmt.Errorf("failed to decode tcp server: %w", err)
			}
			// Устанавливаем значения по умолчанию для TCP сервера
			if s.Port == 0 {
				s.Port = defaultPort
			}
			if s.Host == "" {
				s.Host = defaultHost
			}
			if s.MaxConnections == 0 {
				s.MaxConnections = defaultMaxConnections
			}
			if s.MaxMessageSize == 0 {
				s.MaxMessageSize = defaultMaxMessageSize
			}
			if s.IdleTimeout == 0 {
				s.IdleTimeout = defaultIdleTimeout
			}
			server = &s

		default:
			return fmt.Errorf("unknown server type: %s", baseServer.Type)
		}

		servers = append(servers, server)
	}

	*s = servers
	return nil
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
