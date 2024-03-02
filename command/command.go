package command

import (
	"bot/client"
	"github.com/SevereCloud/vksdk/v2/events"
)

type Command interface {
	Name() string
	Execute(args []string, user *client.User, object events.MessageNewObject)
}
