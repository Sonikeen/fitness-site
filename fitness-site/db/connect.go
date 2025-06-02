package db

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
    var err error
    Pool, err = pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5432/fitness-site")
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }
    fmt.Println("✅ Подключено к базе данных (pgxpool)")
}
