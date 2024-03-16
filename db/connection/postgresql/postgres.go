package postgresql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUsername string
	DbName     string
	DbPassword string
}

func OpenPg() (*sql.DB, error) {
	conf := Config{
		DbHost:     os.Getenv("DB_HOST"),
		DbName:     os.Getenv("DB_NAME"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.DbUsername,
		conf.DbPassword,
		conf.DbHost,
		conf.DbPort,
		conf.DbName,
	)
	if os.Getenv("ENV") == "production" {
		connStr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=verify-full sslrootcert=ap-southeast-1-bundle.pem Timezone=UTC",
			conf.DbUsername,
			conf.DbPassword,
			conf.DbHost,
			conf.DbPort,
			conf.DbName,
		)
	}

	db, err := sql.Open("postgres", connStr)

	return db, err
}
