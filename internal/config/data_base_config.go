package config

type DataBaseConfig struct {
	PgPort   string `env:"PG_PORT" envDefault:"5432"`
	PgHost   string `env:"PG_HOST" envDefault:"localhost"`
	PgDBName string `env:"PG_DB_NAME" envDefault:"chatdb"`
	PgUser   string `env:"PG_USER" envDefault:"postgres"`
	PgPwd    string `env:"PG_PWD" envDefault:"159753"`
}
