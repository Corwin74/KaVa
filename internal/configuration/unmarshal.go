package configuration

import (
	"fmt"

	"code.cloudfoundry.org/bytefmt"
	"gopkg.in/yaml.v3"
)

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