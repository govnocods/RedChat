package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/govnocods/RedChat/internal/auth"
	"github.com/govnocods/RedChat/internal/logger"
	"github.com/govnocods/RedChat/models"
	"github.com/govnocods/RedChat/utils"
)

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.UserService.ValidateUser(req.Username, req.Password); err != nil {
		logger.Warn("Registration failed: validation error",
			"username", req.Username,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.UserService.RegisterUser(req.Username, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "user already exists") {
			logger.Warn("Registration failed: user already exists",
				"username", req.Username,
				"ip", r.RemoteAddr,
			)
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
		logger.WithError(err).
			With("username", req.Username).
			With("ip", r.RemoteAddr).
			Error("Failed to register user")
		http.Error(w, "Не удалось зарегистрировать пользователя", http.StatusInternalServerError)
		return
	}

	logger.Info("User registered successfully",
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

	utils.SetCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
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
		if strings.Contains(err.Error(), "invalid credentials") {
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

	utils.SetCookie(w, token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "token generated", "username": "%s", "token": "%s"}`, user.Username, token)
}

func (h *Handlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Logged out successfully",
		"success": true,
	})
}
