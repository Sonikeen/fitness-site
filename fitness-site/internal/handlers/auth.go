package handler

import (
    "encoding/json"
    "net/http"

    "fitness-site/internal/service"
)

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

// APIRegisterHandler — JSON-регистрация (POST /api/register)
func APIRegisterHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    if err := service.Register(req.Name, req.Email, req.Password); err != nil {
        http.Error(w, err.Error(), http.StatusConflict)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("User registered successfully"))
}

// LoginHandler — JSON-вход (POST /api/login)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    user, err := service.Login(req.Email, req.Password)
    if err != nil {
        http.Error(w, "unauthorized", http.StatusUnauthorized)
        return
    }

    w.Write([]byte("Welcome, " + user.Name + "!"))
}
