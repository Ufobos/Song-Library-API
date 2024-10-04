package handlers

import (
	"encoding/json"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/services"
	"strconv"

	"github.com/gorilla/mux"
)

type SongHandler struct {
	service services.SongService
}

func NewSongHandler(service services.SongService) *SongHandler {
	return &SongHandler{service: service}
}

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song to the library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Add song"
// @Success 201 {string} string "Song added successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Group string `json:"group"`
		Song  string `json:"song"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.AddSong(request.Group, request.Song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Song added successfully"))
}

// GetSongs godoc
// @Summary Get songs with filters and pagination
// @Description Retrieve songs based on filters and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Group name"
// @Param song query string false "Song name"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} models.Song
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs [get]
func (h *SongHandler) GetSongs(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]interface{})
	query := r.URL.Query()

	// Filter by fields
	if group := query.Get("group"); group != "" {
		filters["group_name"] = group
	}
	if song := query.Get("song"); song != "" {
		filters["song_name"] = song
	}

	// Pagination
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	limit := 10
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	songs, err := h.service.GetSongs(filters, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

// GetSongText godoc
// @Summary Get song text with pagination over verses
// @Description Retrieve song text with pagination over verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} map[string]string "text"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get pagination parameters
	query := r.URL.Query()
	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")
	limit := 1
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Invalid offset parameter", http.StatusBadRequest)
			return
		}
	}

	text, err := h.service.GetSongText(id, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"text": text})
}

// UpdateSong godoc
// @Summary Update an existing song
// @Description Update song details
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Update song"
// @Success 200 {string} string "Song updated successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var song models.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	song.ID = id

	err = h.service.UpdateSong(&song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Song updated successfully"))
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song by its ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {string} string "Song deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.service.DeleteSong(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Song deleted successfully"))
}
