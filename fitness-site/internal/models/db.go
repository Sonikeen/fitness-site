package models

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // драйвер для database/sql
)

var DB *sql.DB

// InitDB открывает соединение через database/sql + pgx driver и проверяет Ping.
func InitDB(connString string) error {
	var err error
	DB, err = sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("models.InitDB: не удалось открыть соединение: %w", err)
	}
	if err = DB.PingContext(context.Background()); err != nil {
		return fmt.Errorf("models.InitDB: не удалось опинговать БД: %w", err)
	}
	fmt.Println("✅ models.DB инициализирован")
	return nil
}
