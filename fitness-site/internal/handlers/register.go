package handlers

import (
    "net/http"
)

// RegisterHandlers — если нужен отдельный обработчик 
func RegisterHandlers(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        ShowRegister(w, r)
    case http.MethodPost:
        HandleRegister(w, r)
    default:
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
    }
}
