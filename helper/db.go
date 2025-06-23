package helper

import (
	"database/sql"
	"fmt"
	migrate "khaira-admin/db"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
)

func NewDb() (*sql.DB, func(), error) {
	//WITHOUT DOCKER!
	// err := godotenv.Load()
	// if err != nil {
	// 	return nil, nil, err
	// }

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsnMigrate := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	migrate.RunMigrations(dsnMigrate)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	cleanup := func() {
		db.Close()
	}

	return db, cleanup, nil
}
