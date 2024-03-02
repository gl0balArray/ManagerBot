package command

import (
	"bot/chat"
	"bot/client"
	internalerrors "bot/error"
	api "github.com/SevereCloud/vksdk/v2/api"
	"github.com/sirupsen/logrus"
)

type Pool struct {
	Logger   logrus.Logger
	Api      *api.VK
	UserPool *client.UserPool
	ChatPool *chat.Pool
	Commands map[string]Command
}

func NewCommandPool(logger logrus.Logger, userPool *client.UserPool, chatPool *chat.Pool, vkApi *api.VK) *Pool {
	pool := &Pool{logger, vkApi, userPool, chatPool, map[string]Command{}}

	pool.Commands["old"] = OldCommand{Pool: pool}
	pool.Commands["status"] = StatusCommand{Pool: pool}
	pool.Commands["profile"] = ProfileCommand{Pool: pool}

	//manage
	pool.Commands["staff"] = StaffCommand{Pool: pool}

	return pool
}

func (pool *Pool) Get(name string) (Command, error) {
	if cmd, isset := pool.Commands[name]; isset {
		return cmd, nil
	}

	return EmptyCommand{}, &internalerrors.CommandNotFound{}
}
