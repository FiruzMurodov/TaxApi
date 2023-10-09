package handlers

import (
	"encoding/json"
	"net/http"
	"taxApi/internal/models"
)

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrInvalidData))
		h.Lg.Error()
		return
	}

	err = h.Service.ValidateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrWrongNumberChar))
		h.Lg.Error()
		return
	}

	err = h.Service.IsLoginUsed(user.Login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertToJson(models.ErrLoginUsed))
		h.Lg.Error()
		return
	}

	err = h.Service.RegistrationUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(ConvertToJson(models.ErrInternal))
		h.Lg.Error()
		return
	}

	_, err = w.Write(ConvertToJson(models.SuccessRegistration))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.Lg.Error()
		return
	}
}
