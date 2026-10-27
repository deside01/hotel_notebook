package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/deside01/hotel_notebook/backend/internal/db"
	"github.com/deside01/hotel_notebook/backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type GuestRecordHandler struct {
	repo db.Querier
}

func NewGuestRecordHandler(repo db.Querier) *GuestRecordHandler {
	return &GuestRecordHandler{
		repo: repo,
	}
}

type guestRecordBody struct {
	FullName string `json:"full_name,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Room     string `json:"room,omitempty"`
	CheckIn  string `json:"check_in,omitempty"`
	CheckOut string `json:"check_out,omitempty"`
	Note     string `json:"note,omitempty"`
}

func (h *GuestRecordHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var body guestRecordBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid json body")
		return
	}

	if body.FullName == "" {
		utils.Err(w, http.StatusBadRequest, "full_name field is required")
		return
	}

	if body.Phone == "" {
		utils.Err(w, http.StatusBadRequest, "phone field is required")
		return
	}

	if body.Room == "" {
		utils.Err(w, http.StatusBadRequest, "room field is required")
		return
	}

	if body.CheckIn == "" {
		utils.Err(w, http.StatusBadRequest, "check_in field is required")
		return
	}

	if body.CheckOut == "" {
		utils.Err(w, http.StatusBadRequest, "check_out field is required")
		return
	}

	checkIn, err := time.Parse("2006-01-02", body.CheckIn)
	if err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid check_in field")
		return
	}

	checkOut, err := time.Parse("2006-01-02", body.CheckOut)
	if err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid check_out field")
		return
	}

	record, err := h.repo.CreateGuestRecord(r.Context(), db.CreateGuestRecordParams{
		UserID:   int64(userID),
		FullName: body.FullName,
		Phone:    body.Phone,
		Room:     body.Room,
		CheckIn: pgtype.Date{
			Time:  checkIn,
			Valid: true,
		},
		CheckOut: pgtype.Date{
			Time:  checkOut,
			Valid: true,
		},
		Note: body.Note,
	})

	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, record)
}

func (h *GuestRecordHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	records, err := h.repo.GetGuestRecordsByUserID(r.Context(), int64(userID))
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, records)
}

func (h *GuestRecordHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	record, err := h.repo.GetGuestRecordByID(r.Context(), int64(id))
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userID != int(record.UserID) {
		utils.Err(w, http.StatusForbidden, "forbidden")
		return
	}

	utils.JSON(w, http.StatusOK, record)
}

func (h *GuestRecordHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	var body guestRecordBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid json body")
		return
	}

	record, err := h.repo.GetGuestRecordByID(r.Context(), int64(id))
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userID != int(record.UserID) {
		utils.Err(w, http.StatusForbidden, "forbidden")
		return
	}

	checkOut, err := time.Parse("2006-01-02", body.CheckOut)
	if err != nil {
		utils.Err(w, http.StatusBadRequest, "invalid check_out field")
		return
	}

	err = h.repo.UpdateGuestRecord(r.Context(), db.UpdateGuestRecordParams{
		ID: record.ID,
		CheckOut: pgtype.Date{
			Time:  checkOut,
			Valid: true,
		},
	})
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusNoContent, nil)
}

func (h *GuestRecordHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	record, err := h.repo.GetGuestRecordByID(r.Context(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.Err(w, http.StatusNotFound, "not found")
			return
		}

		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userID != int(record.UserID) {
		utils.Err(w, http.StatusForbidden, "forbidden")
		return
	}

	err = h.repo.DeleteGuestRecordByID(r.Context(), record.ID)
	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusNoContent, nil)
}

func (h *GuestRecordHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 20
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	result, err := h.repo.SearchGuests(r.Context(), db.SearchGuestsParams{
		FullName:  query.Get("full_name"),
		Phone:     query.Get("phone"),
		Room:      query.Get("room"),
		ReqLimit:  int32(limit),
		ReqOffset: int32(offset),
	})

	if err != nil {
		utils.Err(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, result)
}
