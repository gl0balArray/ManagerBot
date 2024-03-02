package listener

import (
	"bot/chat"
	"bot/client"
	"bot/command"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/sirupsen/logrus"
	"strings"
)

type Listener struct {
	Interface   *longpoll.LongPoll
	UserPool    *client.UserPool
	CommandPool *command.Pool
	ChatPool    *chat.Pool
	Logger      logrus.Logger
}

func NewListener(interfaz *longpoll.LongPoll, cmdPool *command.Pool, chatPool *chat.Pool, logger logrus.Logger, userPool *client.UserPool) *Listener {
	return &Listener{interfaz, userPool, cmdPool, chatPool, logger}
}

/*func (h *Handler) OnUpdate(update api.Update) {
	if update.Message != nil {
		h.OnMessage(update.Message)
	}
}*/

func (l *Listener) OnMessage(obj events.MessageNewObject) {

	l.Logger.Info(obj.Message.Action)
	l.Logger.Info(obj.Message.FromID)

	peerId := obj.Message.PeerID
	id := obj.Message.FromID

	user, err := l.UserPool.Get(id)

	if err != nil {
		return
	}

	actionType := obj.Message.Action.Type
	memberId := obj.Message.Action.MemberID

	if actionType != "" && peerId != id { //conversation
		ch, err := l.ChatPool.Get(peerId)

		if err != nil {
			return
		}

		target, err := l.UserPool.Get(memberId)

		if err != nil {
			return
		}

		if actionType == "chat_kick_user" {
			ch.OnKick(user, target)
			return
		} else if actionType == "chat_invite_user" {
			ch.OnInvite(user, target)
			return
		}
	}

	args := strings.Fields(obj.Message.Text)
	if len(args) < 1 {
		return
	}

	cmdData := strings.Split(args[0], "/")

	if cmdData[0] == "" {
		if cmd, err := l.CommandPool.Get(cmdData[1]); err == nil {
			if err != nil {
				l.Logger.Error(err)
				return
			}

			go cmd.Execute(args, user, obj)
		}
	}
}
