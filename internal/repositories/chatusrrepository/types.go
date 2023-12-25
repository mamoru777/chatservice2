package chatusrrepository

import "github.com/google/uuid"

type ChatUsr struct {
	ChatId uuid.UUID `db:"chat_id"`
	UsrId  uuid.UUID `db:"usr_id"`
}
