package chat

import (
	"bot/client"
	"bot/storage"
	"bot/storage/template"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"log"
	"sync"
)

type Chat struct {
	Api        *api.VK
	Data       *template.ChatData
	Provider   *storage.MySQLProvider
	MembersMux *sync.Mutex
	Members    map[int]*Member
}

func (chat *Chat) LoadOtherData() {
	chat.LoadStaff()
}

func (chat *Chat) LoadStaff() {
	q := fmt.Sprintf("select * from vk_chat_admins where peer_id = %d", chat.Data.PeerID)

	staff := struct {
		Admins []*template.ChatAdmin
	}{}

	rows, err := chat.Provider.Driver().NamedQuery(q, &staff)

	if err != nil {
		return
	}

	for rows.Next() {
		if err := rows.StructScan(staff); err != nil {
			log.Println(err)
			return
		}
	}

	for _, admin := range staff.Admins {
		chat.Members[admin.ID].Priority = admin.Priority
	}
}

func (chat *Chat) SaveAdmin(member *Member) {
	memberId := member.User.Data.VkId
	if member.Priority == 0 {
		q := fmt.Sprintf("delete from vk_chat_admins where peer_id = %d and user_id = %d", chat.Data.PeerID, memberId)

		_, err := chat.Provider.Driver().Query(q)

		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	q := fmt.Sprintf("insert into vk_chat_admins (peer_id, user_id, priority) VALUES(%d, %d, %d)"+
		"ON DUPLICATE KEY UPDATE priority = %d", chat.Data.PeerID, memberId, member.Priority, member.Priority,
	)

	_, err := chat.Provider.Driver().Query(q)

	if err != nil {
		log.Println(err.Error())
	}
}

func (chat *Chat) OnInvite(host, target *client.User) {
	chat.MembersMux.Lock()
	chat.Members[target.Data.VkId] = &Member{
		0,
		target,
		true,
	}
	chat.MembersMux.Unlock()

	chat.LoadStaff()
}

func (chat *Chat) OnKick(admin, target *client.User) {
	targetId := target.Data.VkId
	chat.MembersMux.Lock()
	member, exists := chat.Members[targetId]
	chat.MembersMux.Unlock()

	prms := params.NewMessagesSendBuilder()
	prms.PeerID(chat.Data.PeerID)
	prms.DisableMentions(true)
	prms.RandomID(0)

	if admin.Data.VkId == targetId { //self
		prms.Message(fmt.Sprintf("%s вышел из чата.", target.Mention()))
		_, _ = chat.Api.MessagesSend(prms.Params)
		return
	}

	prms.Message(fmt.Sprintf("%s исключил %s из чата.", admin.Mention(), target.Mention()))
	_, _ = chat.Api.MessagesSend(prms.Params)

	if !exists {
		chat.SaveAdmin(&Member{0, target, true})
		return
	}

	if member.Priority != 0 {
		member.Priority = 0
		chat.SaveAdmin(member)
	}

	delete(chat.Members, targetId)
}
