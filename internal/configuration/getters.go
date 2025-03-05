package configuration

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

func (t *TCPServerConfig) getType() string {
	return t.Type
}

func (t *TCPServerConfig) getName() string {
	return t.Name
}