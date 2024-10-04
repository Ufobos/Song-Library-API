package services

import (
	"song-library/configs"
	"song-library/internal/models"
	"song-library/internal/repository"
	"song-library/internal/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

type SongService interface {
	AddSong(group, song string) error
	GetSongs(filters map[string]interface{}, limit, offset int) ([]models.Song, error)
	GetSongText(id string, limit, offset int) (string, error)
	UpdateSong(song *models.Song) error
	DeleteSong(id string) error
}

type songService struct {
	repo   repository.SongRepository
	config *configs.Config
}

func NewSongService(repo repository.SongRepository, config *configs.Config) SongService {
	return &songService{repo: repo, config: config}
}

func (s *songService) AddSong(group, song string) error {
	newSong := models.Song{
		GroupName:   group,
		SongName:    song,
		ReleaseDate: time.Now(),
		Text:        "Тестовый текст песни.",
		Link:        "https://example.com",
	}

	err := s.repo.Create(&newSong)
	if err != nil {
		utils.Logger.Error("Не удалось сохранить песню в базе данных", zap.Error(err))
		return err
	}

	utils.Logger.Info("Песня успешно добавлена", zap.String("group", group), zap.String("song", song))
	return nil
}

// GetSongs retrieves songs based on filters and pagination.
func (s *songService) GetSongs(filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	return s.repo.FilterWithPagination(filters, limit, offset)
}

// GetSongText retrieves the text of a song with pagination over verses.
func (s *songService) GetSongText(id string, limit, offset int) (string, error) {
	song, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}

	verses := strings.Split(song.Text, "\n\n")
	start := offset
	end := offset + limit
	if start >= len(verses) {
		return "", nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return strings.Join(verses[start:end], "\n\n"), nil
}

// UpdateSong updates an existing song.
func (s *songService) UpdateSong(song *models.Song) error {
	return s.repo.Update(song)
}

// DeleteSong deletes a song by its ID.
func (s *songService) DeleteSong(id string) error {
	return s.repo.Delete(id)
}
