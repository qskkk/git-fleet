package command

var Handled = map[string]func(group string) (out string, err error){
	"status": ExecuteStatus,
	"ls":     ExecuteStatus,
}

var GlobalHandled = map[string]func(group string) (out string, err error){
	"config": ExecuteConfig,
	"help":   ExecuteHelp,
	"status": ExecuteStatus,
	"ls":     ExecuteStatus,
}
