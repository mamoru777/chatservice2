package messagerepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type MessageRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (mr *MessageRepository) Create(ctx context.Context, m *Message) error {
	const q = `
		INSERT INTO messages (chat_id, usr_id, text, data) 
			VALUES (:chat_id, :usr_id, :text, :data)
	`
	_, err := mr.db.NamedExecContext(ctx, q, m)
	return err
}

func (mr *MessageRepository) GetList(ctx context.Context, chatId uuid.UUID) ([]*Message, error) {
	const q = `
		SELECT * FROM messages WHERE chat_id = $1
	`
	var list []*Message
	err := mr.db.SelectContext(ctx, &list, q, chatId)
	return list, err
}
