package role

import "errors"

type Pool struct {
	Roles map[string]Role
}

func NewPool() *Pool {
	pool := &Pool{map[string]Role{}}

	pool.Roles[USER] = Role{USER, 0, "Пользователь"}
	pool.Roles[OWNER] = Role{OWNER, 100, "Владелец"}

	return pool
}

func (pool *Pool) Get(id string) (Role, error) {
	if role, contains := pool.Roles[id]; contains {
		return role, nil
	}

	return pool.Roles[USER], errors.New("role not found")
}
