package db

import (
	"context"
	"errors"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	SqlConn *sqlx.DB
	RdbConn *redis.Client
}

func SqlConnect() (*sqlx.DB, error) {
	// conn, err := sqlx.Connect("pgx", "postgresql://postgres.pgnjdntiwyriphrlnxhp:hAdwad-3naqgo-fewtup@aws-0-eu-central-1.pooler.supabase.com:6543/postgres")
	// conn, err := sqlx.Connect("pgx", "postgresql://postgres.pgnjdntiwyriphrlnxhp:hAdwad-3naqgo-fewtup@aws-0-eu-central-1.pooler.supabase.com:5432/postgres")
	// conn, err := sqlx.Connect("pgx", "user=postgres.pgnjdntiwyriphrlnxhp password=hAdwad-3naqgo-fewtup host=aws-0-eu-central-1.pooler.supabase.com port=6543 dbname=postgres")
	conn, err := sqlx.Connect("pgx", "user=postgres.pgnjdntiwyriphrlnxhp password=hAdwad-3naqgo-fewtup host=aws-0-eu-central-1.pooler.supabase.com port=5432 dbname=postgres")
	log.Print("Connected to PostgreSQL")
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return conn, nil
}

func RedisConnect(ctx context.Context) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis",
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return rdb, nil
}

func Connect(ctx context.Context) (*DB, error) {
	postgres, err := SqlConnect()
	if err != nil {
		return nil, errors.New("Failed to connect to PostgreSQL")
	}

	return &DB{
		SqlConn: postgres,
	}, nil
}

func (db *DB) Close() (err error) {
	if sqlErr := db.SqlConn.Close(); sqlErr != nil {
		log.Printf("Failed to close PostgreSQL connection: %v", sqlErr)
		return sqlErr
	}
	log.Print("Closed connection to PostgreSQL")

	if rdbErr := db.RdbConn.Close(); rdbErr != nil {
		log.Printf("Failed to close Redis connection: %v", rdbErr)
		return rdbErr
	}

	return nil
}
