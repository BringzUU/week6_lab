package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"moviehub/mw"
	"moviehub/store"
)

// ทั้งสามเส้นทางสวมประตู login และใช้เลขสมาชิกจากตั๋วเท่านั้น
// จึงไม่มีทางขอรายการอยากดูของคนอื่นได้

func (h *Handler) ListWishlist(w http.ResponseWriter, r *http.Request) {
	who := mw.WhoFrom(r.Context())
	items, err := h.Store.ListWishlist(r.Context(), who.ID)
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

type wishlistReq struct {
	Note *string `json:"note"`
}

func (h *Handler) PutWishlist(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "movieId")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	var req wishlistReq
	// ไม่ส่งเนื้อคำขอมาเลยก็ได้ (EOF) ถือว่าไม่มีโน้ต
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	if req.Note != nil {
		n := strings.TrimSpace(*req.Note)
		if utf8.RuneCountInString(n) > 200 {
			writeError(w, http.StatusBadRequest, "NOTE_TOO_LONG", "โน้ตยาวได้ไม่เกิน 200 ตัวอักษร")
			return
		}
		if n == "" {
			req.Note = nil
		} else {
			req.Note = &n
		}
	}
	who := mw.WhoFrom(r.Context())
	err := h.Store.PutWishlist(r.Context(), who.ID, movieID, req.Note)
	if errors.Is(err, store.ErrNoMovie) {
		writeError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "ไม่พบหนังเรื่องนี้")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"movieId": movieID, "note": req.Note})
}

func (h *Handler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	movieID, ok := pathID(r, "movieId")
	if !ok {
		writeError(w, http.StatusBadRequest, "BAD_ID", "เลขหนังต้องเป็นตัวเลข")
		return
	}
	who := mw.WhoFrom(r.Context())
	if err := h.Store.DeleteWishlist(r.Context(), who.ID, movieID); err != nil {
		internal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
