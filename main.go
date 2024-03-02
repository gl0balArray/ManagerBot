package main

import (
	"bot/chat"
	"bot/client"
	"bot/command"
	"bot/listener"
	"bot/storage"
	"bot/util"
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	api "github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/sirupsen/logrus"
)

var logger logrus.Logger
var interfaz *longpoll.LongPoll
var vkApi *api.VK
var DataProvider *storage.MySQLProvider

var CommandPool *command.Pool
var UserPool *client.UserPool
var ChatPool *chat.Pool

var groupId int = 194456885

func main() {
	logger = util.NewLogger()

	logger.Info("Setting up")

	os.Setenv("VK_TOKEN", "token")

	vkApi = api.NewVK(os.Getenv("VK_TOKEN"))

	DataProvider = &storage.MySQLProvider{
		Username: "username",
		Password: "pwd",
		Address:  "localhost",
		Port:     "3306",
		Database: "db",
		Logger:   logger,
	}

	DataProvider.Setup()

	UserPool = client.NewUserPool(DataProvider, logger, vkApi)

	ChatPool = chat.NewPool(DataProvider, vkApi, UserPool, logger)

	CommandPool = command.NewCommandPool(logger, UserPool, ChatPool, vkApi)

	lp, err := longpoll.NewLongPoll(vkApi, groupId)

	if err != nil {
		logger.Fatal(err)
	}

	interfaz = lp

	handler := listener.NewListener(interfaz, CommandPool, ChatPool, logger, UserPool)

	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		handler.OnMessage(obj)
	})

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for now := range time.Tick(3 * time.Minute) {
			for _, user := range UserPool.Users {
				q := "UPDATE `vk_userdata` SET `nickname` = :nickname, `role` = :role WHERE id = " + strconv.Itoa(user.Data.ID)
				if _, err := DataProvider.Driver().NamedExec(q, user.Data); err != nil {
					logger.Fatal(err)
				}
			}

			logger.Info("[" + now.String() + "] Data is saved")
		}
	}()

	go func() {
		<-signalChan
		logger.Warning("Shutting down")
		os.Exit(1)
	}()

	if err := lp.Run(); err != nil {
		logger.Fatal(err)
	}
}
