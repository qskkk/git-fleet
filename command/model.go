package command

var Handled = map[string]func() (out string, err error){
	"status": ExecuteStatus,
	"ls":     ExecuteStatus,
}

var GlobalHandled = map[string]func() (out string, err error){
	"config": ExecuteStatus,
	"help":   ExecuteStatus,
	"status": ExecuteStatus,
}
