package client

import (
	"bot/client/role"
	internalerrors "bot/error"
	"bot/storage"
	"bot/storage/template"
	"errors"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type UserPool struct {
	UsersMux sync.Mutex
	Api      *api.VK
	Users    map[int]*User
	RolePool *role.Pool
	Provider *storage.MySQLProvider
	Logger   logrus.Logger
}

func NewUserPool(provider *storage.MySQLProvider, logger logrus.Logger, vkApi *api.VK) *UserPool {
	return &UserPool{sync.Mutex{}, vkApi, map[int]*User{}, role.NewPool(), provider, logger}
}

func (pool *UserPool) Get(id int) (*User, error) {
	if user, contains := pool.Users[id]; contains {
		return user, nil
	}

	user := &User{Pool: pool}

	data, err := pool.LoadData(id)

	if err != nil {
		var dataNotFoundError internalerrors.DataNotFoundError
		if errors.As(err, &dataNotFoundError) {
			return pool.Get(id)
		}
		return nil, err
	}

	user.Data = data

	user.LoadData()

	pool.UsersMux.Lock()

	pool.Users[id] = user

	pool.UsersMux.Unlock()

	return user, nil
}

func (pool *UserPool) LoadData(id int) (*template.ClientData, error) {

	pool.Logger.Info("Trying to get data for id" + strconv.Itoa(id))

	q := "select * from `vk_userdata` where vk_id = " + strconv.Itoa(id)

	data := &template.ClientData{}

	rows, err := pool.Provider.Driver().NamedQuery(q, data)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.StructScan(&data); err != nil {
			return nil, err
		}
	}

	defer rows.Close()

	if data.VkId != id {
		pool.Logger.Info("Registering new user (userid: " + strconv.Itoa(id) + ")")
		q = "INSERT INTO `vk_userdata`(`vk_id`) VALUES(:vk_id)"
		if _, err := pool.Provider.Driver().NamedExec(q, map[string]interface{}{
			"vk_id": id,
		}); err != nil {
			return nil, err
		}
		return nil, internalerrors.DataNotFoundError{
			What: "not found",
		}
	}

	return data, nil
}
