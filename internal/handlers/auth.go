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
	"github.com/govnocods/RedChat/internal/logger"
	"github.com/govnocods/RedChat/models"
)

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.RegisterUser(req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			logger.Warn("Registration failed: user already exists",
				"username", req.Username,
				"ip", r.RemoteAddr,
			)
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		logger.WithError(err).
			With("username", req.Username).
			With("ip", r.RemoteAddr).
			Error("Failed to register user")
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	logger.Info("User registered successfully",
		"user_id", user.Id,
		"username", user.Username,
		"ip", r.RemoteAddr,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "User registered successfully",
		"username": user.Username,
		"id":       user.Id,
	})
}

func (h *Handlers) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), "invalid credentials") {
			logger.Warn("Authentication failed",
				"username", req.Username,
				"ip", r.RemoteAddr,
				"reason", "invalid credentials",
			)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		logger.WithError(err).
			With("username", req.Username).
			With("ip", r.RemoteAddr).
			Error("Database error during authentication")
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	logger.Info("User authenticated successfully",
		"user_id", user.Id,
		"username", user.Username,
		"ip", r.RemoteAddr,
	)

	token, err := auth.GenerateToken(user.Id, user.Username)
	if err != nil {
		logger.WithError(err).
			With("user_id", user.Id).
			With("username", user.Username).
			Error("Failed to generate token")
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
