package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deside01/hotel_notebook/backend/internal/db"
	"github.com/deside01/hotel_notebook/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo db.Querier
}

func NewAuthHandler(repo db.Querier) *AuthHandler {
	return &AuthHandler{
		repo: repo,
	}
}

type authBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var body authBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if body.Email == "" {
		utils.Err(w, http.StatusBadRequest, "email field is required")
		return
	}

	if body.Password == "" {
		utils.Err(w, http.StatusBadRequest, "password field is required")
		return
	}

	fetchedUser, err := h.repo.FindUserByEmail(r.Context(), body.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.Err(w, http.StatusInternalServerError, "internal server error")
		return
	}
	if fetchedUser.Email != "" {
		utils.Err(w, http.StatusConflict, "user already exists")
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	user, err := h.repo.CreateUser(r.Context(), db.CreateUserParams{
		Email:    body.Email,
		Password: string(pass),
	})
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.JSON(w, http.StatusCreated, user)
}
