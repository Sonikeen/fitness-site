// internal/models/program.go
package models

import (
	"context"
	"errors"
)

// Program описывает фитнес-программу.
type Program struct {
	ID          int
	Name        string
	Description string
	Days        []string // Список описаний тренировок: индекс 0 = День 1, и т. д.
}

var ErrProgramNotFound = errors.New("program not found")

// В этой версии мы загружаем программы через ProgramStorage, поэтому
// GetAllPrograms и GetProgramByID из этого файла не используются напрямую.

func GetAllPrograms(ctx context.Context) ([]Program, error) {
	return nil, nil
}

func GetProgramByID(ctx context.Context, programID int) (*Program, error) {
	return nil, ErrProgramNotFound
}
