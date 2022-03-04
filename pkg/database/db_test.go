package database

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestConnectDB(t *testing.T) {
	conf := &Config{
		Username: "postgres",
		Password: "postgres",
		Database: "postgres",
		Address:  "localhost",
		Port:     5432,
	}
	_, err := ConnectDB(conf)
	if err != nil {
		t.Error("Error connecting to database")
	}
	// defer db.Close()
}
