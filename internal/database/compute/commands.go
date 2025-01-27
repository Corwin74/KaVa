package compute


const (
	UnknownCommandID = iota
	SetCommandID
	GetCommandID
	DelCommandID
)

const (
	setCommand = "SET"
	getCommand = "GET"
	delCommand = "DEL"
)

var commandTextToID = map[string]int{
	setCommand: SetCommandID,
	getCommand: GetCommandID,
	delCommand: DelCommandID,
}


var commandArgumentsCount = map[int]int{
	SetCommandID: 2,
	GetCommandID: 1,
	DelCommandID: 1,
}