package chatrepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (cr *ChatRepository) Create(ctx context.Context) (uuid.UUID, error) {
	const q = `
		INSERT INTO chats DEFAULT VALUES RETURNING id
	`
	var id uuid.UUID
	err := cr.db.QueryRowContext(ctx, q).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, err
}

func (cr *ChatRepository) Get(ctx context.Context, chatId uuid.UUID) (*Chat, error) {
	const q = `
		SELECT id FROM chats WHERE id = $1
	`
	s := new(Chat)
	err := cr.db.GetContext(ctx, s, q, chatId)
	return s, err
}

func (cr *ChatRepository) GetList(ctx context.Context) ([]*Chat, error) {
	const q = `
		SELECT * FROM chats 
	`
	var list []*Chat
	err := cr.db.SelectContext(ctx, &list, q)
	return list, err
}
