package command

import (
	"bot/client"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"runtime"
	"strconv"
)

type StatusCommand struct {
	Pool *Pool
	Command
}

func (cmd StatusCommand) Name() string {
	return "status"
}

func (cmd StatusCommand) Execute(args []string, user *client.User, object events.MessageNewObject) {
	prms := params.NewMessagesSendBuilder()
	prms.PeerID(object.Message.PeerID)
	prms.RandomID(0)

	prms.Message(
		"GoroutineNum: " + strconv.Itoa(runtime.NumGoroutine()) + "\n" +
			"NumCPU: " + strconv.Itoa(runtime.NumCPU()) + "\n" +
			"Version: " + runtime.Version())

	_, _ = cmd.Pool.Api.MessagesSend(prms.Params)

	cmd.Pool.Logger.Info(user.Data.Role)
}
