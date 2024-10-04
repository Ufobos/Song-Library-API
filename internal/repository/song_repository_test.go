package repository

import (
	"song-library/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateSong(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewSongRepository(sqlxDB)

	song := &models.Song{
		GroupName:   "Muse",
		SongName:    "Supermassive Black Hole",
		ReleaseDate: time.Now(),
		Text:        "Sample text",
		Link:        "http://example.com",
	}

	mock.ExpectPrepare("INSERT INTO songs").
		ExpectExec().
		WithArgs(song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Link).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(song)
	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
