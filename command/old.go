package command

import (
	"bot/client"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"math/rand"
)

type OldCommand struct {
	Pool *Pool
	Command
}

func (cmd OldCommand) Name() string {
	return "old"
}

func (cmd OldCommand) Execute(args []string, user *client.User, object events.MessageNewObject) {
	prms := params.NewMessagesSendBuilder()
	prms.PeerID(object.Message.PeerID)
	prms.RandomID(0)

	startMessageId := 0

	if rand.Intn(1) >= 5 {
		startMessageId = object.Message.ConversationMessageID - rand.Intn(3)
	} else {
		startMessageId = object.Message.ConversationMessageID + rand.Intn(3)
	}

	messages := []int{}

	for i := 0; i < 99; i++ {
		messages = append(messages, startMessageId+i)
	}

	prms.ForwardMessages(messages)
	prms.Message("wtf")

	_, _ = cmd.Pool.Api.MessagesSend(prms.Params)
}
