package usrrepository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsrRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UsrRepository {
	return &UsrRepository{
		db: db,
	}
}

func (cr *UsrRepository) Create(ctx context.Context, c *Usr) error {
	const q = `
		INSERT INTO usrs (id) 
			VALUES (:id)
	`
	_, err := cr.db.NamedExecContext(ctx, q, c)
	return err
}

func (cr *UsrRepository) Get(ctx context.Context, usrId uuid.UUID) (*Usr, error) {
	const q = `
		SELECT id FROM usrs WHERE id = $1
	`
	s := new(Usr)
	err := cr.db.GetContext(ctx, s, q, usrId)
	return s, err
}

func (cr *UsrRepository) GetList(ctx context.Context) ([]*Usr, error) {
	const q = `
		SELECT * FROM usrs 
	`
	var list []*Usr
	err := cr.db.SelectContext(ctx, &list, q)
	return list, err
}
