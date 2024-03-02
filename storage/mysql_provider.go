package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type MySQLProvider struct {
	Username string
	Password string
	Database string
	Address  string
	Port     string
	Logger   log.Logger
	DB       *sqlx.DB
}

func (provider *MySQLProvider) Setup() {
	dsn := provider.Username + ":" + provider.Password + "@tcp(" + provider.Address + ":" + provider.Port + ")/" + provider.Database + "?parseTime=true"
	driver, err := sqlx.Open("mysql", dsn)
	if err != nil {
		provider.Logger.Fatal("Couldn't connect to DB")
		return
	}
	driver.SetMaxOpenConns(10)
	driver.SetMaxIdleConns(10)
	provider.SetDriver(driver)
}

func (provider *MySQLProvider) Ping() {
	driver := provider.Driver()
	if err := driver.Ping(); err != nil {
		provider.Logger.Info("not return ping, creating a new connection")
		provider.Setup()
	}
}

func (provider *MySQLProvider) Driver() *sqlx.DB {
	return provider.DB
}

func (provider *MySQLProvider) SetDriver(db *sqlx.DB) {
	provider.DB = db
}
