package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/biryanim/SongLibrary/internal/entities"
	"github.com/biryanim/SongLibrary/pkg/errors"
	"github.com/biryanim/SongLibrary/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"strings"
)

type SongServiceUseCase interface {
	GetSongs(ctx context.Context, group, name, page, limit string) ([]entities.Song, error)
	GetSongById(ctx context.Context, id int) (*entities.Song, error)
	DeleteSongById(ctx context.Context, id int) error
	UpdateSongById(ctx context.Context, id int, song entities.Song) error
	PostSong(ctx context.Context, song *entities.Song) error
}

type Adapter struct {
	s SongServiceUseCase
}

func New(s SongServiceUseCase) *Adapter {
	return &Adapter{
		s: s,
	}
}

func StartServer(port int, a *Adapter) {
	route := chi.NewRouter()

	route.Use(logger.RequestLogger)

	route.Get("/songs", a.SongPage)
	route.Get("/songs/{id}/lyrics", a.SongLyrics)
	route.Delete("/songs/{id}", a.DeleteSong)
	route.Put("/songs/{id}", a.UpdateSong)
	route.Post("/songs", a.PostSong)

	logger.Log.Info("Listening on port", zap.Int("port", port))
	logger.Log.Fatal("failed to start server", zap.Error(http.ListenAndServe(fmt.Sprintf(":%d", port), route)))
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), route))
}

// @Summary list
// @Tags songs
// @Description song representation
// @Accept json
// @Produce json
// @Param  group   query string  false  "name search by group"
// @Param  name   query string  false  "song name search by name"
// @Param  page   query int  false  "page"
// @Param  limit   query int  false  "numbers of songs per page"
// @Success 200 {object} entities.Song
// @Failure 500 {object} errors.ErrorResponse
// @Router       /songs [get]
func (a *Adapter) SongPage(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	songName := chi.URLParam(r, "name")
	page := chi.URLParam(r, "page")
	limit := chi.URLParam(r, "limit")

	songs, err := a.s.GetSongs(r.Context(), group, songName, page, limit)

	if err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Debug("", zap.Any("url", r.URL.String()), zap.Any("numbers of songs", len(songs)), zap.Any("songs", songs))

	w.Header().Add("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(songs); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// @Summary song
// @Tags lyrics
// @Description song lyrics representation
// @Accept json
// @Produce json
// @Param  id   path int  true  "song id"
// @Param  verse   query int  false  "verse of the song"
// @Param  limit   query int  true  "count of verses"
// @Success 200 {object} entities.Song
// @Failure 400	{object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router       /songs/{id}/lyrics [get]
func (a *Adapter) SongLyrics(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid id parameter")
		return
	}

	verse, err := strconv.Atoi(r.URL.Query().Get("verse"))
	if err != nil || verse < 1 {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid verse parameter")
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 1
	}

	song, err := a.s.GetSongById(r.Context(), id)

	if err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	paginatedLyrics := splitLyrics(song.Lyrics, verse, limit)

	logger.Log.Debug("", zap.Any("url", r.URL.String()), zap.Any("lyrics", paginatedLyrics))

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(paginatedLyrics); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

func splitLyrics(lyrics string, verse, limit int) []string {
	verses := strings.Split(lyrics, "\n\n")
	startInd := (verse - 1) * limit
	endInd := startInd + limit
	if endInd > len(verses) {
		endInd = len(verses)
	}
	return verses[startInd:endInd]
}

// @Summary delete
// @Tags song
// @Description delete song
// @Accept json
// @Produce json
// @Param  id path int true "song id"
// @Success 200 {object} Response
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router       /cars/{id} [delete]
func (a *Adapter) DeleteSong(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid id parameter")
		return
	}

	err = a.s.DeleteSongById(r.Context(), id)

	if err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Debug("Deleted", zap.Any("song", id))

	response := Response{
		Message: "Song successfully deleted",
		ID:      id,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// @Summary update
// @Tags song
// @Description update song data by ID
// @Accept json
// @Produce json
// @Param  id path int true "song id"
// @Param  input body   entities.Song   true  "song struct"
// @Success 200 {object} Response
// @Failure 400 {object} errors.ErrorResponse
// @Failure 500 {object} errors.ErrorResponse
// @Router       /songs/{id} [put]
func (a *Adapter) UpdateSong(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid id parameter")
		return
	}

	song, err := a.s.GetSongById(r.Context(), id)

	var input entities.Song
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.GroupName != "" {
		song.GroupName = input.GroupName
	}
	if input.SongName != "" {
		song.SongName = input.SongName
	}
	if input.Lyrics != "" {
		song.Lyrics = input.Lyrics
	}
	if input.ReleaseDate != "" {
		song.ReleaseDate = input.ReleaseDate
	}
	if input.Link != "" {
		song.Link = input.Link
	}

	if err = a.s.UpdateSongById(r.Context(), id, *song); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Debug("Updated", zap.Any("song", song))

	response := Response{
		Message: "Song successfully updated",
		ID:      id,
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}

// @Summary add
// @Tags  song
// @Description add song
// @Accept json
// @Produce json
// @Param  song body entities.Song  true  "song name and group name"
// @Success 200 {object} entities.Song
// @Failure 400	{object} errors.ErrorResponse
// @Failure 500	{object} errors.ErrorResponse
// @Router       /songs [post]
func (a *Adapter) PostSong(w http.ResponseWriter, r *http.Request) {
	var song entities.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	songDetail, statusCode, err := externalAPIRequest(song.GroupName, song.SongName)
	if err != nil {
		errors.ServerErrorResponse(w, r, statusCode, err.Error())
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Lyrics = songDetail.Lyrics
	song.Link = songDetail.Link

	if err = a.s.PostSong(r.Context(), &song); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Log.Debug("Posted", zap.Any("song", song))

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(songDetail); err != nil {
		errors.ServerErrorResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}
}
