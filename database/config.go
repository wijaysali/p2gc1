package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
}

func NewConfig() *Config {
	return &Config{
		Username: "root",
		Password: "",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "p2gc1",
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port, c.DBName)
}

func InitDB() (*sql.DB, error) {
	config := NewConfig()
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
