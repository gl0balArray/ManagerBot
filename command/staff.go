package command

import (
	ch "bot/chat"
	"bot/client"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"log"
	"sort"
	"strings"
)

type StaffCommand struct {
	Command
	Pool *Pool
}

func (cmd StaffCommand) Name() string {
	return "staff"
}

func (cmd StaffCommand) Execute(args []string, user *client.User, object events.MessageNewObject) {
	peerId := object.Message.PeerID

	messagesBuilder := params.NewMessagesSendBuilder()
	messagesBuilder.RandomID(0)
	messagesBuilder.PeerID(peerId)
	messagesBuilder.DisableMentions(true)

	if user.Data.VkId == peerId {
		messagesBuilder.Message(fmt.Sprintf("%s, данная команда доступна только в беседах!", user.Mention()))
		_, _ = cmd.Pool.Api.MessagesSend(messagesBuilder.Params)
		return
	}

	chat, err := cmd.Pool.ChatPool.Get(peerId)

	if err != nil {
		messagesBuilder.Message(err.Error())
		_, _ = cmd.Pool.Api.MessagesSend(messagesBuilder.Params)
		return
	}

	staff := []*ch.Member{}

	message := user.Mention() + ", Список персонала данной беседы:\n\n"

	log.Println(chat.Members[577945732])

	for _, member := range chat.Members {
		role, err := cmd.Pool.UserPool.RolePool.Get(member.User.Data.Role)

		if err != nil {
			if role.Priority >= member.Priority {
				staff = append(staff, member)
				continue
			}
		}

		if member.Priority > 0 {
			staff = append(staff, member)
		}
	}

	sort.SliceStable(staff, func(i, j int) bool {
		return staff[j].Priority <= staff[i].Priority
	})

	lastPriorityRank := ch.CreatorPriority

	for _, admin := range staff {
		if admin.Priority == ch.CreatorPriority {
			message += fmt.Sprintf("\nСоздатель чата: ")
		}

		if admin.Priority < lastPriorityRank {
			lastPriorityRank = admin.Priority
			message = strings.TrimSuffix(message, ", ")

			rank, err := cmd.Pool.ChatPool.StaffPool.Get(admin.Priority)

			if err != nil {
				//тут ищи бля
				continue
			}
			message += "\n\n" + rank.NameTag + ": "
		}

		if admin.Priority == ch.CreatorPriority {
			message += "👑 "
		} else if !admin.CanKick {
			message += "⭐ "
		}

		message += fmt.Sprintf("[id%d|%s], ", admin.User.Data.VkId, admin.User.FullName)
	}

	message = strings.TrimSuffix(message, ", ")

	messagesBuilder.Message(message)

	_, _ = cmd.Pool.Api.MessagesSend(messagesBuilder.Params)
}
