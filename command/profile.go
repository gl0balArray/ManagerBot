package command

import (
	"bot/client"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
)

type ProfileCommand struct {
	Pool *Pool
	Command
}

func (cmd ProfileCommand) Name() string {
	return "profile"
}

func (cmd ProfileCommand) Execute(args []string, user *client.User, object events.MessageNewObject) {
	prms := params.NewMessagesSendBuilder()

	prms.RandomID(0)
	prms.PeerID(object.Message.PeerID)
	prms.DisableMentions(true)

	cmd.Pool.Logger.Info(object.Message.PeerID)

	role, _ := cmd.Pool.UserPool.RolePool.Get(user.Data.Role)

	profileMessage := `
%s, Информация о вашем профиле:

🎫 WhNiggaID: %d
🎖 Роль: %s
📅 Дата регистрации: %s
`

	prms.Message(fmt.Sprintf(profileMessage, user.Mention(), user.Data.ID, role.FormattedName, user.Data.RegisterTimestamp.Format("02 January 2006")))

	_, _ = cmd.Pool.Api.MessagesSend(prms.Params)
}
