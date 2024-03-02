package chat

import (
	"bot/client"
	"bot/storage"
	"bot/storage/template"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type Pool struct {
	ChatsMux  *sync.Mutex
	Logger    logrus.Logger
	Api       *api.VK
	UserPool  *client.UserPool
	StaffPool *StaffPool
	Provider  *storage.MySQLProvider
	Chats     map[int]*Chat
}

func NewPool(provider *storage.MySQLProvider, vkApi *api.VK, userPool *client.UserPool, logger logrus.Logger) *Pool {
	staffPool := NewStaffPool()
	return &Pool{&sync.Mutex{}, logger, vkApi, userPool, staffPool, provider, map[int]*Chat{}}
}

func (pool *Pool) Get(id int) (*Chat, error) {
	if chat, contains := pool.Chats[id]; contains {
		return chat, nil
	}

	chat, err := pool.LoadData(id)

	if err != nil {
		return nil, err
	}

	chat.LoadOtherData()

	pool.ChatsMux.Lock()
	pool.Chats[id] = chat
	pool.ChatsMux.Unlock()

	return chat, nil
}

func (pool *Pool) LoadData(id int) (*Chat, error) {
	pool.Logger.Info("Trying to get data for chat" + strconv.Itoa(id))

	q := fmt.Sprintf("select * from `vk_chatdata` where peer_id = %d", id)

	data := &template.ChatData{}

	rows, err := pool.Provider.Driver().NamedQuery(q, &data)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.StructScan(&data); err != nil {
			return nil, err
		}
	}

	chat := &Chat{MembersMux: &sync.Mutex{}, Provider: pool.Provider, Api: pool.Api}
	chat.MembersMux.Lock()
	chat.Members = map[int]*Member{}

	defer chat.MembersMux.Unlock()

	membersBuilder := params.NewMessagesGetConversationMembersBuilder()
	membersBuilder.PeerID(id)

	response, err := pool.Api.MessagesGetConversationMembers(membersBuilder.Params)

	if err != nil {
		return nil, err
	}

	for _, member := range response.Items {
		user, err := pool.UserPool.Get(member.MemberID)

		if err != nil || member.MemberID < 0 {
			continue
		}

		chatMember := &Member{
			User:     user,
			CanKick:  member.CanKick,
			Priority: 0,
		}

		if member.IsOwner {
			chatMember.Priority = CreatorPriority
		} else if member.IsAdmin {
			chatMember.Priority = AdminPriority
		}

		chat.Members[member.MemberID] = chatMember
	}

	if data.PeerID != id {
		pool.Logger.Info("Registering new chat (peer_id: " + strconv.Itoa(id) + ")")

		q = "INSERT INTO `vk_chatdata`(`peer_id`) VALUES(:peer_id)"

		if _, err := pool.Provider.Driver().NamedExec(q, map[string]interface{}{
			"peer_id": id,
		}); err != nil {
			return nil, err
		}
	}

	chat.Data = data

	return chat, nil
}
