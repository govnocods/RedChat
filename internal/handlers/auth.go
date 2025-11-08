package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/govnocods/RedChat/internal/auth"
	"github.com/govnocods/RedChat/models"
)

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.UserService.RegisterUser(req.Username, req.Password)
}

func (h *Handlers) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		// Логируем для себя
		fmt.Println("Auth error:", err)

		// Обрабатываем конкретно
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "invalid credentials") {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	/*
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
	*/

	token, err := auth.GenerateToken(user.Id, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "token generated", "username": "%s", "token": "%s"}`, user.Username, token)
}
