package handlers

import (
	"encoding/json"
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

type regBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Registration(w http.ResponseWriter, r *http.Request) {
	var body regBody
	json.NewDecoder(r.Body).Decode(&body)

	if body.Email == "" {
		utils.Err(w, 400, "email field is required")
		return
	}

	if body.Password == "" {
		utils.Err(w, 400, "password field is required")
		return
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)

	user, _ := h.repo.CreateUser(r.Context(), db.CreateUserParams{
		Email:    body.Email,
		Password: string(pass),
	})

	utils.JSON(w, 201, user)
}
