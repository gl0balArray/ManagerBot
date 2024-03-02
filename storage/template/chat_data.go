package template

import "time"

type ChatData struct {
	ID                int       `db:"id"`
	PeerID            int       `db:"peer_id"`
	IsGold            bool      `db:"gold"`
	RegisterTimestamp time.Time `db:"created_at"`
}
