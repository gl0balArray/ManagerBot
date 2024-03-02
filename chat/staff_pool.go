package chat

import "errors"

const (
	AdminPriority   = 90
	CreatorPriority = 100
)

type StaffPool struct {
	Ranks map[int]*Rank
}

func NewStaffPool() *StaffPool {
	pool := &StaffPool{map[int]*Rank{}}

	pool.Ranks[AdminPriority] = &Rank{"Администратор", AdminPriority}
	pool.Ranks[CreatorPriority] = &Rank{"Создатель", CreatorPriority}

	return pool
}

func (pool *StaffPool) Get(priority int) (*Rank, error) {
	rank, contains := pool.Ranks[priority]

	if !contains {
		return nil, errors.New("rank not found")
	}

	return rank, nil
}
