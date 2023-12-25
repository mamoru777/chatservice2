package chatusrrepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChatUsrRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *ChatUsrRepository {
	return &ChatUsrRepository{
		db: db,
	}
}

func (cur *ChatUsrRepository) Create(ctx context.Context, cu *ChatUsr) error {
	const q = `
		INSERT INTO chat_usr (chat_id, usr_id) 
			VALUES (:chat_id, :usr_id)
	`
	_, err := cur.db.NamedExecContext(ctx, q, cu)
	return err
}

func (cur *ChatUsrRepository) Get(ctx context.Context, usrId uuid.UUID, friendId uuid.UUID) (uuid.UUID, error) {
	const query = `
		SELECT chat_id
		FROM chat_usr
		WHERE usr_id IN ($1, $2)
		GROUP BY chat_id
		HAVING COUNT(DISTINCT usr_id) = 2
	`
	var chatID uuid.UUID
	err := cur.db.QueryRow(query, usrId, friendId).Scan(&chatID)
	if err != nil {
		return uuid.Nil, err
	}

	return chatID, nil
}
