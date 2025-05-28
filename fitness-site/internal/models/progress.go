package models

import "time"

// Progress хранит отметку о выполнении дня пользователем
type Progress struct {
    ID          int
    UserID      int
    ProgramID   int
    Day         int       // номер дня
    CompletedAt time.Time
}
