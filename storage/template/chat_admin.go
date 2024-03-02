package template

type ChatAdmin struct {
	ID       int `db:"user_id"`
	Priority int `db:"priority"`
}
