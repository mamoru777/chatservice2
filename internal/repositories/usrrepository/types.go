package usrrepository

import "github.com/google/uuid"

type Usr struct {
	Id uuid.UUID `db:"id"`
}
