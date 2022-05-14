package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"githb.com/demo-employee-api/internal/config"
)

func NewDb(conf *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Name)
	fmt.Println(dsn)

	connDb, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = connDb.Ping()

	if err != nil {
		return nil, err
	}

	return connDb, nil
}

// type DB interface {
// 	Create()
// }

// type DB struct {
// 	db *sql.DB
// }
