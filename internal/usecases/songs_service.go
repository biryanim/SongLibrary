package usecases

import (
	"context"
	"github.com/biryanim/SongLibrary/internal/entities"
	"strconv"
)

type storage interface {
	GetALlSongs(ctx context.Context, group, name string, page, limit int) ([]entities.Song, error)
	GetSong(ctx context.Context, id int) (*entities.Song, error)
	DeleteSong(ctx context.Context, id int) error
	UpdateSong(ctx context.Context, id int, song entities.Song) error
	CreateSong(ctx context.Context, song entities.Song) error
}

type SongsService struct {
	storage storage
}

func New(st storage) *SongsService {
	return &SongsService{
		storage: st,
	}
}

func (s *SongsService) GetSongs(ctx context.Context, group, name, p, l string) ([]entities.Song, error) {
	limit, _ := strconv.Atoi(p)
	offset, _ := strconv.Atoi(l)
	if limit == 0 {
		limit = 1
	}
	if offset == 0 {
		offset = 10
	}

	return s.storage.GetALlSongs(ctx, "%"+group+"%", "%"+name+"%", limit, (limit-1)*offset)
}

func (s *SongsService) GetSongById(ctx context.Context, id int) (*entities.Song, error) {
	return s.storage.GetSong(ctx, id)
}

func (s *SongsService) DeleteSongById(ctx context.Context, id int) error {
	return s.storage.DeleteSong(ctx, id)
}

func (s *SongsService) UpdateSongById(ctx context.Context, id int, song entities.Song) error {

	return s.storage.UpdateSong(ctx, id, song)
}

func (s *SongsService) PostSong(ctx context.Context, song *entities.Song) error {
	return s.storage.CreateSong(ctx, *song)
}
