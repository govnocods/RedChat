package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/govnocods/RedChat/internal/auth"
	"github.com/govnocods/RedChat/models"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if h.DB == nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	req.Password = string(hashedPassword)

	h.DB.CreateUser(req)
}

func (h *Handlers) AuthHAndler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.DB.GetUser(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "token generated", "username": "%s", "token": "%s"}`, user.Username, token)
}
