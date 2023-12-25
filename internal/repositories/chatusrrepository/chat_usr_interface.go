package chatusrrepository

import (
	"context"
	"github.com/google/uuid"
)

type IChatUsrRepository interface {
	Create(ctx context.Context, cu *ChatUsr) error
	Get(ctx context.Context, usrId uuid.UUID, friendId uuid.UUID) (uuid.UUID, error)
	//Update(ctx context.Context, c *Chat) error
	//Delete(ctx context.Context, chatId uuid.UUID) error
}
