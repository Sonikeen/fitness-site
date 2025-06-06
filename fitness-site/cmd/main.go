package main

import (
    "log"
    "net/http"

    "fitness-site/db"
    "fitness-site/internal/handlers"
    "fitness-site/internal/service"
    "fitness-site/internal/storage"
    "fitness-site/internal/router"
)

func main() {
    // 1) Подключаемся к базе — внутри инициализируется models.DB и создаётся pgxpool.Pool
    db.Connect()

    // 2) Создаём storage (используем глобальный пул db.Pool)
    userStore    := storage.NewUserStorage(db.Pool)
    programStore := storage.NewProgramStorage(db.Pool)
    progressStore := storage.NewProgressStorage(db.Pool)
    workoutStore  := storage.NewWorkoutStorage()

    // 3) Инициализируем сервисы и присваиваем глобальным переменным в handlers
    handlers.UserService     = service.NewUserService(userStore)
    handlers.ProgramService  = service.NewProgramService(programStore)
    handlers.ProgressService = service.NewProgressService(progressStore)
    handlers.WorkoutService  = service.NewWorkoutService(workoutStore)

    // 4) Запускаем готовый router
    r := router.SetupRouter()

    log.Println("Сервер слушает :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatalf("Ошибка запуска сервера: %v", err)
    }
}
