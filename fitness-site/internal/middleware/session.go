package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type sessionData struct {
	UserID    int
	ExpiresAt time.Time
}

var (
	globalSessionStore = map[string]sessionData{}
	mu                 sync.Mutex
	sessionTTL         = 7 * 24 * time.Hour
)

func generateSID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CreateSession(userID int) (string, error) {
	sid, err := generateSID()
	if err != nil {
		return "", err
	}
	mu.Lock()
	globalSessionStore[sid] = sessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionTTL),
	}
	mu.Unlock()
	return sid, nil
}

func GetUserIDBySession(sid string) (int, error) {
	mu.Lock()
	data, ok := globalSessionStore[sid]
	mu.Unlock()
	if !ok {
		return 0, errors.New("session not found")
	}
	if time.Now().After(data.ExpiresAt) {
		mu.Lock()
		delete(globalSessionStore, sid)
		mu.Unlock()
		return 0, errors.New("session expired")
	}
	return data.UserID, nil
}

func DeleteSession(sid string) {
	mu.Lock()
	delete(globalSessionStore, sid)
	mu.Unlock()
}
