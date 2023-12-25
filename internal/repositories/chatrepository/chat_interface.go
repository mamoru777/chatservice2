package chatrepository

import (
	"context"
	"github.com/google/uuid"
)

type IChatRepository interface {
	Create(ctx context.Context) (uuid.UUID, error)
	Get(ctx context.Context, chatId uuid.UUID) (*Chat, error)
	GetList(ctx context.Context) ([]*Chat, error)
	//Update(ctx context.Context, c *Chat) error
	//Delete(ctx context.Context, chatId uuid.UUID) error
}
