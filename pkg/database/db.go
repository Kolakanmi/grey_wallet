package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Kolakanmi/grey_wallet/pkg/envconfig"
)

//Config - db config struct
type Config struct {
	Username string `envconfig:"DATABASE_USERNAME" default:"postgres"`
	Password string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	Database string `envconfig:"DATABASE_NAME" default:"postgres"`
	Address  string `envconfig:"DATABASE_ADDRESS" default:"localhost"`
	Port     int    `envconfig:"DATABASE_PORT" default:"5432"`
}

var tables = []string{
	`CREATE TABLE kola_transactions (
		id VARCHAR(50) PRIMARY KEY NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
		updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
		deleted_at timestamp with time zone DEFAULT NULL,
		amount decimal(10,2) NOT NULL,
		status text NOT NULL DEFAULT 'pending'
	)`,
	`CREATE TABLE kola_wallets (
		id VARCHAR(50) PRIMARY KEY NOT NULL,
		created_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
		updated_at timestamp with time zone NOT NULL DEFAULT current_timestamp,
		deleted_at timestamp with time zone DEFAULT NULL,
		balance decimal(10,2) NOT NULL
	)`,
}

func LoadConfig() *Config {
	var conf Config
	envconfig.Load(&conf)
	return &conf
}

func ConnectDB(conf *Config) (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Address, conf.Port, conf.Username, conf.Password, conf.Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	loadTables(db)

	return db, nil

}

func loadTables(db *sql.DB) error {
	for _, table := range tables {
		_, err := db.Exec(table)
		if err != nil {
			return err
		}
	}
	return nil
}
