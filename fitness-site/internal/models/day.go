package models

// Day описывает один день программы
type Day struct {
    DayNumber   int
    Description string
    ImageURL    string // URL картинки (может быть пустым)
}
