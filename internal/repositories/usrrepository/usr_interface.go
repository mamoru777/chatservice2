package usrrepository

import (
	"context"
	"github.com/google/uuid"
)

type IUsrRepository interface {
	Create(ctx context.Context, u *Usr) error
	Get(ctx context.Context, usrId uuid.UUID) (*Usr, error)
	GetList(ctx context.Context) ([]*Usr, error)
	//Update(ctx context.Context, c *Usr) error
	//Delete(ctx context.Context, usrId uuid.UUID) error
}
