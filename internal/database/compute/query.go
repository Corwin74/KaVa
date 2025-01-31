package compute

// Query - распарсенный запрос
type Query struct {
	commandID int
	key       string
	value     string
}

// NewQuery -- конструктор
func NewQuery(commandID int, key, value string) Query {
	return Query{
		commandID: commandID,
		key:       key,
		value:     value,
	}
}

// CommandID -- getter
func (q *Query) CommandID() int {
	return q.commandID
}

// GetKey -- getter
func (q *Query) GetKey() string {
	return q.key
}

// GetValue -- getter
func (q *Query) GetValue() string {
	return q.value
}
