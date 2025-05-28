package handlers

import (
    "html/template"
    "net/http"
    "path/filepath"
    "strconv"

    "github.com/go-chi/chi/v5"

    "fitness-site/db"
    "fitness-site/internal/middleware"
    "fitness-site/internal/service"
    "fitness-site/internal/storage"
)

var (
    programService  *service.ProgramService
    progressService *service.ProgressService
)

func init() {
    db.Connect()
    progStore := storage.NewProgramStorage(db.Conn)
    programService = service.NewProgramService(progStore)

    progProgStore := storage.NewProgressStorage(db.Conn)
    progressService = service.NewProgressService(progProgStore)
}

// ProgramList — защищённый список программ (GET /programs/)
func ProgramList(w http.ResponseWriter, r *http.Request) {
    progs, err := programService.GetAllPrograms()
    if err != nil {
        http.Error(w, "Ошибка получения программ", http.StatusInternalServerError)
        return
    }
    render(w, "programs.html", map[string]interface{}{
        "Programs": progs,
    })
}

// ProgramDetail — детали одной программы и прогресс (GET /programs/{id})
func ProgramDetail(w http.ResponseWriter, r *http.Request) {
    // 1) ID программы из URL
    pid, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.NotFound(w, r)
        return
    }
    // 2) Достаем программу
    prog, err := programService.GetProgramByID(pid)
    if err != nil {
        http.Error(w, "Программа не найдена", http.StatusNotFound)
        return
    }

    // 3) Получаем userID из сессии/контекста
    userIDVal := r.Context().Value(middleware.UserIDKey)
    if userIDVal == nil {
        http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
        return
    }
    userID := userIDVal.(int)

    // 4) Загружаем прогресс
    progresses, err := progressService.ListProgress(userID, pid)
    if err != nil {
        http.Error(w, "Не удалось загрузить прогресс", http.StatusInternalServerError)
        return
    }
    completed := make(map[int]bool, len(progresses))
    for _, p := range progresses {
        completed[p.Day] = true
    }

    // 5) Рендерим шаблон
    data := map[string]interface{}{
        "Program":       prog,
        "CompletedDays": completed,
    }
    t := template.Must(template.ParseFiles(
        filepath.Join("internal", "templates", "base.html"),
        filepath.Join("internal", "templates", "program_detail.html"),
    ))
    t.ExecuteTemplate(w, "base", data)
}
