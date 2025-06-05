package handlers
import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"fitness-site/internal/middleware"
	"fitness-site/internal/models"
)
type DashboardItem struct {
	Program   models.Program
	Completed int
	Total     int
}
type UserStats struct {
	TotalCompletedDays int
	TotalStarted       int
	TotalFinished      int
	AvgProgress        int 
}
type DashboardData struct {
	ActiveTab  string
	Error      string
	Items      []DashboardItem
	Stats      *UserStats
	IsLoggedIn bool
	IsAdmin    bool
}
func Dashboard(w http.ResponseWriter, r *http.Request) {
	basePath := filepath.Join("internal", "templates", "base.html")
	dashPath := filepath.Join("internal", "templates", "dashboard.html")
	t, err := template.ParseFiles(basePath, dashPath)
	if err != nil {
		log.Printf("Dashboard.Handler: ошибка при ParseFiles: %v", err)
		http.Error(w, "Ошибка серверного шаблона", http.StatusInternalServerError)
		return
	}
	user, ok := middleware.GetUser(r)
	if !ok {
		active := "login"
		if r.URL.Query().Get("tab") == "register" {
			active = "register"
		}
		data := DashboardData{
			ActiveTab:  active,
			Error:      "",
			Items:      nil,
			Stats:      nil,
			IsLoggedIn: false,
			IsAdmin:    false,
		}
		t.ExecuteTemplate(w, "base", data)
		return
	}
	userID := user.ID
	progs, err := ProgramService.GetAllPrograms(r.Context())
	if err != nil {
		http.Error(w, "Ошибка получения программ", http.StatusInternalServerError)
		return
	}
	var items []DashboardItem
	for _, p := range progs {
		doneList, err := ProgressService.ListProgress(r.Context(), userID, p.ID)
		if err != nil {
			http.Error(w, "Ошибка загрузки прогресса", http.StatusInternalServerError)
			return
		}
		doneSet := make(map[int]struct{})
		for _, pr := range doneList {
			doneSet[pr.DayNumber] = struct{}{}
		}
		items = append(items, DashboardItem{
			Program:   p,
			Completed: len(doneSet),
			Total:     len(p.Days),
		})
	}
	var statDays, started, finished, sumPercent int
	for _, it := range items {
		if it.Completed > 0 {
			started++
		}
		if it.Total > 0 && it.Completed == it.Total {
			finished++
		}
		statDays += it.Completed
		if it.Total > 0 {
			sumPercent += it.Completed * 100 / it.Total
		}
	}
	avg := 0
	if len(items) > 0 {
		avg = sumPercent / len(items)
	}
	stats := &UserStats{
		TotalCompletedDays: statDays,
		TotalStarted:       started,
		TotalFinished:      finished,
		AvgProgress:        avg,
	}
	data := DashboardData{
		ActiveTab:  "",
		Error:      "",
		Items:      items,
		Stats:      stats,
		IsLoggedIn: true,
		IsAdmin:    user.IsAdmin, // передаём, является ли пользователь админом
	}
	t.ExecuteTemplate(w, "base", data)
}
