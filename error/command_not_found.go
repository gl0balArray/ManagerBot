package error

type CommandNotFound struct{

}

func (e *CommandNotFound) Error() string{
	return ""
}