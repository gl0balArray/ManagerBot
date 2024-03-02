package role

const (
	USER  = "user"
	OWNER = "owner"
)

type Role struct {
	ID            string
	Priority      int
	FormattedName string
}

func (r Role) HasRights(priority int) bool {
	return r.Priority >= priority
}
