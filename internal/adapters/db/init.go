package db

import (
	"context"
	"fmt"
	"github.com/biryanim/SongLibrary/internal/entities"
	"github.com/jackc/pgx/v5"
)

type Adapter struct {
	db *pgx.Conn
}

func New(conn *pgx.Conn) *Adapter {
	return &Adapter{
		db: conn,
	}
}

func (a *Adapter) GetALlSongs(ctx context.Context, group, name string, limit, offset int) ([]entities.Song, error) {
	query := `
	SELECT  id, group_name, song_name, lyrics, release_date, link
	FROM songs
	WHERE group_name LIKE $1 AND song_name LIKE $2
	LIMIT $3 OFFSET $4`
	rows, err := a.db.Query(ctx, query, group, name, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var songs []entities.Song
	for rows.Next() {
		s := entities.Song{}
		err = rows.Scan(&s.ID, &s.SongName, &s.GroupName, &s.SongName, &s.Lyrics, &s.ReleaseDate, &s.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, s)
	}
	return songs, nil
}

func (a *Adapter) GetSong(ctx context.Context, id int) (*entities.Song, error) {
	query := `
	SELECT  id, group_name, song_name, lyrics, release_date, link
	FROM songs
	WHERE id = $1
	`

	row := a.db.QueryRow(ctx, query, id)
	s := entities.Song{}
	err := row.Scan(&s.ID, &s.SongName, &s.GroupName, &s.SongName, &s.Lyrics, &s.ReleaseDate, &s.Link)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (a *Adapter) DeleteSong(ctx context.Context, id int) error {
	query := `
	DELETE FROM songs WHERE id = $1`
	result, err := a.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no song found with ID %d", id)
	}
	return nil
}

func (a *Adapter) UpdateSong(ctx context.Context, id int, song entities.Song) error {
	query := `
	UPDATE songs
	SET group_name = $1, song_name =$2, lyrics = $3, release_date = $4, link = $5
	WHERE id = $6`

	_, err := a.db.Exec(ctx, query, song.GroupName, song.SongName, song.Lyrics, song.ReleaseDate, song.Link)

	return err
}

func (a *Adapter) CreateSong(ctx context.Context, song entities.Song) error {
	query := `
	INSERT INTO songs (group_name, song_name, lyrics, release_date, link)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := a.db.Exec(ctx, query, song.GroupName, song.SongName, song.Lyrics, song.ReleaseDate, song.Link)

	return err
}
