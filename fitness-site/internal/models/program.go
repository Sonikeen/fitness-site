package models

// Program описывает тренировочную программу
type Program struct {
    ID          int
    Name        string
    Description string
    Days        []string // описания каждого дня
}
