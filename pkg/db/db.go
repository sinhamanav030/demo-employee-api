package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"githb.com/demo-employee-api/internal/config"
)

func NewDb(conf *config.Config, logger *log.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", conf.Database.User, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.Name)
	// fmt.Println(dsn)

	connDb, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	err = connDb.Ping()

	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return connDb, nil
}
