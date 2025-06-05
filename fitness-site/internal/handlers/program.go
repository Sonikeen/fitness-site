package handlers
import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"github.com/go-chi/chi/v5"
	"fitness-site/internal/middleware"
	"fitness-site/internal/models"
)
type ProgramDetailPageData struct {
	Program       *models.Program // ID, Name, Description, Days []string
	CompletedDays map[int]bool    // номера выполненных дней
}
func ProgramPageHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/dashboard?tab=login", http.StatusSeeOther)
		return
	}
	idStr := chi.URLParam(r, "id")
	progID, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	prog, err := ProgramService.GetProgramByID(r.Context(), progID)
	if err != nil || prog == nil {
		http.NotFound(w, r)
		return
	}
	doneList, err := ProgressService.ListProgress(r.Context(), userID, progID)
	if err != nil {
		http.Error(w, "Ошибка загрузки прогресса", http.StatusInternalServerError)
		return
	}
	completedMap := make(map[int]bool, len(doneList))
	for _, pr := range doneList {
		completedMap[pr.DayNumber] = true
	}
	t, err := template.New("").
		Funcs(template.FuncMap{"add": func(a, b int) int { return a + b }}).
		ParseFiles(
			filepath.Join("internal", "templates", "base.html"),
			filepath.Join("internal", "templates", "program_detail.html"),
		)
	if err != nil {
		log.Printf("ProgramPageHandler: ошибка при парсинге шаблонов: %v", err)
		http.Error(w, "Ошибка серверного шаблона", http.StatusInternalServerError)
		return
	}
	data := ProgramDetailPageData{
		Program:       prog,
		CompletedDays: completedMap,
	}
	if execErr := t.ExecuteTemplate(w, "base", data); execErr != nil {
		log.Printf("ProgramPageHandler: ExecuteTemplate вернул ошибку: %v", execErr)
		http.Error(w, "Ошибка серверного рендеринга: "+execErr.Error(), http.StatusInternalServerError)
		return
	}
}
type workoutInfoJSON struct {
	DayNumber int    `json:"day_number"`
	Date      string `json:"date"`
	Exercises string `json:"exercises"`
	Notes     string `json:"notes"`
}
type programDetailJSON struct {
	ID          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Workouts    []workoutInfoJSON `json:"workouts"`
}
func ProgramHandlerJSON(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	programID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid program ID", http.StatusBadRequest)
		return
	}
	if _, ok := middleware.UserIDFromContext(r.Context()); !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	prog, err := ProgramService.GetProgramByID(r.Context(), programID)
	if err != nil || prog == nil {
		http.Error(w, "Программа не найдена", http.StatusNotFound)
		return
	}
	workouts, err := models.GetWorkoutsByProgram(r.Context(), programID)
	if err != nil {
		http.Error(w, "Ошибка при получении тренировок", http.StatusInternalServerError)
		return
	}
	var resp programDetailJSON
	resp.ID = prog.ID
	resp.Title = prog.Name
	resp.Description = prog.Description
	for _, wkt := range workouts {
		resp.Workouts = append(resp.Workouts, workoutInfoJSON{
			DayNumber: wkt.DayNumber,
			Date:      wkt.Date.Format("2006-01-02"),
			Exercises: wkt.Exercises,
			Notes:     wkt.Notes,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}