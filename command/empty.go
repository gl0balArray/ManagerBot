package command

import (
	"bot/client"
	"github.com/SevereCloud/vksdk/v2/events"
)

type EmptyCommand struct {
	Command
}

func (cmd EmptyCommand) Name() string {
	return "empty"
}

func (cmd EmptyCommand) Execute(args []string, user *client.User, object events.MessageNewObject) {
}
