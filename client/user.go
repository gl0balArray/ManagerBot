package client

import (
	"bot/storage/template"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"strconv"
)

type User struct {
	Data         *template.ClientData
	Pool         *UserPool
	CacheMention string
	FullName     string
}

func (u *User) LoadData() {
	if u.Data.VkId < 0 { //TODO.
		u.CacheMention = "group"
		u.FullName = "group"
		return
	}

	mention := "[id" + strconv.Itoa(u.Data.VkId) + "|"

	if u.Data.Nickname != "NaN" {
		mention += u.Data.Nickname + "]"

		u.CacheMention = mention
	}

	usersGetBuilder := params.NewUsersGetBuilder()
	usersGetBuilder.UserIDs([]string{"id" + strconv.Itoa(u.Data.VkId)})

	if response, err := u.Pool.Api.UsersGet(usersGetBuilder.Params); err == nil {
		user := response[0]

		u.FullName = user.FirstName + " " + user.LastName

		if u.CacheMention == "" {
			mention += u.FullName + "]"

			u.CacheMention = mention
		}
	}
}

func (u *User) Mention() string {
	if u.CacheMention != "" {
		return u.CacheMention
	}

	return ""
}
