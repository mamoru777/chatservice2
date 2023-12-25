package chatrepository

import "github.com/google/uuid"

type Chat struct {
	Id uuid.UUID `db:"id"`
}
