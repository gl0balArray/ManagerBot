package chat

import (
	"bot/client"
	"github.com/SevereCloud/vksdk/v2/object"
)

type Member struct {
	Priority int
	User     *client.User
	CanKick  object.BaseBoolInt
}
