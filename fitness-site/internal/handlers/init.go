package handlers

import (
    "fitness-site/internal/service"
)

// Глобальные сервисы (назначаются в main.go)
var (
    UserService     *service.UserService
    ProgramService  *service.ProgramService
    ProgressService *service.ProgressService
    WorkoutService  *service.WorkoutService
)
