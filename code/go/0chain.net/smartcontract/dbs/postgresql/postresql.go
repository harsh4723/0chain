package postgresql

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	"0chain.net/smartcontract/dbs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetPostgresSqlDb(config dbs.DbAccess) (dbs.Store, error) {
	if !config.Enabled {
		return nil, nil
	}
	db := &PostgresStore{}
	err := db.Open(config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type PostgresStore struct {
	db *gorm.DB
}

func (store *PostgresStore) Open(config dbs.DbAccess) error {
	if !config.Enabled {
		return errors.New("db_open_error, db disabled")
	}

	maxRetries := 60 * 2 // 2 minutes
	for i := 0; i < maxRetries; i++ {

		db, err := gorm.Open(postgres.Open(fmt.Sprintf(
			"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
			config.Host,
			config.Port,
			config.User,
			config.Name,
			config.Password)),
			&gorm.Config{
				Logger:                 logger.Default.LogMode(logger.Silent),
				SkipDefaultTransaction: true,
				PrepareStmt:            true,
			})

		if err == nil { // tcp host/port are ready
			sqldb, err := db.DB()
			if err == nil {
				err = sqldb.Ping()

				if err == nil { // login/passwd and schema are initialized
					sqldb.SetMaxIdleConns(config.MaxIdleConns)
					sqldb.SetMaxOpenConns(config.MaxOpenConns)
					sqldb.SetConnMaxLifetime(config.ConnMaxLifetime)
					store.db = db
					break
				}
			}
		}

		if (i + 1) == maxRetries {
			return fmt.Errorf("db_open_error, Error opening the DB connection: %v", err)
		}
		fmt.Printf("db: [%v/%v]waiting for postgres to ready\n", i, maxRetries)
		time.Sleep(1 * time.Second)
		continue
	}

	fmt.Println("made event sql database ok")
	return nil
}

func (store *PostgresStore) AutoMigrate() error {
	panic("should not be called")
}

func (store *PostgresStore) Close() {
	if store.db != nil {
		if sqldb, _ := store.db.DB(); sqldb != nil {
			sqldb.Close()
		}
	}
}

func (store *PostgresStore) Get() *gorm.DB {
	return store.db
}
