package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"sync"
)

var instance *sqlx.DB
var isConnExists bool
var mutex sync.Mutex

func GetDB() (*sqlx.DB, error) {
	if isConnExists {
		return instance, nil
	}
	mutex.Lock()
	db, err := sqlx.Connect("postgres", viper.GetString("db.connection"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db. Error: %s", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping db. Error: %s", err)
	}
	isConnExists = true
	instance = db
	mutex.Unlock()
	return instance, nil
}
