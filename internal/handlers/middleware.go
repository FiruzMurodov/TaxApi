package handlers

import (
	"context"
	"net/http"
	"taxApi/internal/models"
)

type ctxKey string

const keyUserID ctxKey = "userId"

func (h *Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		id, err := h.Service.IdByToken(r.Context(), token)
		if err != nil {
			h.Lg.Error()
			w.WriteHeader(http.StatusBadRequest)
			w.Write(ConvertToJson(models.ErrTokenExpired))
			return
		}
		ctx := context.WithValue(r.Context(), keyUserID, id)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
