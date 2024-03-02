package template

import (
	"time"
)

type ClientData struct{
	ID                int       `db:"id"`
	VkId              int       `db:"vk_id"`
	Role              string    `db:"role"`
	Nickname          string    `db:"nickname"`
	RegisterTimestamp time.Time `db:"created_at"`
}