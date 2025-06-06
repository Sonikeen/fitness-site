package db

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func init() {
    // 1) Загружаем .env-файл
    if err := godotenv.Load(); err != nil {
        log.Fatalf("Ошибка при загрузке файла .env: %v", err)
    }

    // 2) Читаем переменные окружения
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    if dbHost == "" || dbPort == "" || dbUser == "" || dbPassword == "" || dbName == "" {
        log.Fatal("Не заданы обязательные переменные окружения для подключения к БД")
    }

    // 3) Формируем строку подключения
    connectionString := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName,
    )

    // 4) Открываем пул соединений
    var err error
    Pool, err = pgxpool.New(context.Background(), connectionString)
    if err != nil {
        log.Fatalf("Ошибка подключения к БД: %v", err)
    }

    fmt.Println("✅ Подключено к базе данных (pgxpool)")
}


func Connect() {
    if Pool == nil {
        log.Fatal("Pool ещё не инициализирован. Возможно, не удалось загрузить .env или подключиться к БД")
    }

}
