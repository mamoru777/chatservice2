package messagerepository

import (
	"github.com/google/uuid"
	"time"
)

type Message struct {
	Id     uuid.UUID `db:"id"`
	ChatID uuid.UUID `db:"chat_id"`
	UsrID  uuid.UUID `db:"usr_id"`
	Text   string    `db:"text"`
	Data   time.Time `db:"data"`
}
