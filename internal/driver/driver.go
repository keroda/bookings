package driver

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

var dbConn = &DB{}

// const maxOpenDbConn = 10
// const maxIdleDbConn = 5
const maxDbConns = 10

const maxDbLifetime = 5 * time.Minute

// urlExample := "postgres://username:password@localhost:5432/database_name"
// win: "postgres://postgres:kr180360@localhost:5432/test_db"
const MyDb = "postgres://kjetilrodal:@localhost:5432/bookings"

// create database pool for Postgres
func ConnectSQL(dsn string) (*pgxpool.Pool, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.Config().MaxConns = maxDbConns
	d.Config().MaxConnLifetime = maxDbLifetime

	dbConn.Pool = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbConn.Pool, nil
}

func testDB(p *pgxpool.Pool) error {
	err := p.Ping(context.Background())
	if err != nil {
		return (err)
	}
	return nil
}

// NewDatabase(dsn string) (*sql.DB, error)
func NewDatabase(dsn string) (*pgxpool.Pool, error) {

	dbpool, err := pgxpool.New(context.Background(), dsn) //os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	//defer dbpool.Close()

	if err = dbpool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return dbpool, nil
}
