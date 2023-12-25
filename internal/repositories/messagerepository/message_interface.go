package messagerepository

import (
	"context"
	"github.com/google/uuid"
)

type IMessageRepository interface {
	Create(ctx context.Context, m *Message) error
	GetList(ctx context.Context, chatId uuid.UUID) ([]*Message, error)
	//Get(ctx context.Context, messageId uuid.UUID) (*Message, error)
	//Update(ctx context.Context, m *Message) error
	//Delete(ctx context.Context, messageId uuid.UUID) error
}
