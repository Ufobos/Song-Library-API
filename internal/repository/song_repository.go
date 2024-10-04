package repository

import (
	"fmt"
	"song-library/internal/models"

	"github.com/jmoiron/sqlx"
)

type SongRepository interface {
	Create(song *models.Song) error
	GetByID(id string) (*models.Song, error)
	Update(song *models.Song) error
	Delete(id string) error
	FilterWithPagination(filters map[string]interface{}, limit, offset int) ([]models.Song, error)
}

type songRepository struct {
	db *sqlx.DB
}

func NewSongRepository(db *sqlx.DB) SongRepository {
	return &songRepository{db: db}
}

// Create adds a new song to the database.
func (r *songRepository) Create(song *models.Song) error {
	query := `INSERT INTO songs (group_name, song_name, release_date, text, link)
              VALUES (:group_name, :song_name, :release_date, :text, :link)`
	_, err := r.db.NamedExec(query, song)
	return err
}

// GetByID retrieves a song by its ID.
func (r *songRepository) GetByID(id string) (*models.Song, error) {
	var song models.Song
	err := r.db.Get(&song, "SELECT * FROM songs WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

// Update modifies an existing song in the database.
func (r *songRepository) Update(song *models.Song) error {
	query := `UPDATE songs SET group_name = :group_name, song_name = :song_name, release_date = :release_date, text = :text, link = :link WHERE id = :id`
	_, err := r.db.NamedExec(query, song)
	return err
}

// Delete removes a song from the database by its ID.
func (r *songRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM songs WHERE id = $1", id)
	return err
}

// FilterWithPagination retrieves songs with filtering and pagination.
func (r *songRepository) FilterWithPagination(filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	query := `SELECT * FROM songs WHERE 1=1`
	args := []interface{}{}
	i := 1

	for k, v := range filters {
		query += fmt.Sprintf(" AND %s = $%d", k, i)
		args = append(args, v)
		i++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	var songs []models.Song
	err := r.db.Select(&songs, query, args...)
	return songs, err
}
