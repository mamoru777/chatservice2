package databaseinit

import (
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mamoru777/chatservice2/internal/config"
)

func InitSqlxDB(cfg config.DataBaseConfig) (*sqlx.DB, error) {
	return sqlx.Connect("pgx", formatConnect(cfg))
}

func formatConnect(cfg config.DataBaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PgUser, cfg.PgPwd, cfg.PgHost, cfg.PgPort, cfg.PgDBName,
	)
}
