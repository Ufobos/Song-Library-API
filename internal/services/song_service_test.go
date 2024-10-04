package services

import (
	"net/http"
	"net/http/httptest"
	"song-library/configs"
	"song-library/internal/models"
	"song-library/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository для имитации методов репозитория
type MockRepository struct {
	mock.Mock
}

// Проверяем, что MockRepository реализует интерфейс SongRepository
var _ repository.SongRepository = (*MockRepository)(nil)

// Реализация методов интерфейса SongRepository

func (m *MockRepository) Create(song *models.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id string) (*models.Song, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Song), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockRepository) Update(song *models.Song) error {
	args := m.Called(song)
	return args.Error(0)
}

func (m *MockRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) FilterWithPagination(filters map[string]interface{}, limit, offset int) ([]models.Song, error) {
	args := m.Called(filters, limit, offset)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Song), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestAddSong(t *testing.T) {
	mockRepo := new(MockRepository)
	config := &configs.Config{
		APIURL: "http://example.com",
	}
	service := NewSongService(mockRepo, config)

	// Мокаем внешний API
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
            "releaseDate": "16.07.2006",
            "text": "Sample text",
            "link": "http://example.com"
        }`))
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Обновляем URL API в конфигурации на адрес тестового сервера
	service.(*songService).config.APIURL = server.URL

	// Настраиваем ожидания
	mockRepo.On("Create", mock.AnythingOfType("*models.Song")).Return(nil)

	err := service.AddSong("Muse", "Supermassive Black Hole")
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetSongText(t *testing.T) {
	mockRepo := new(MockRepository)
	config := &configs.Config{}
	service := NewSongService(mockRepo, config)

	songText := "Verse 1\n\nVerse 2\n\nVerse 3"

	song := &models.Song{
		ID:   1,
		Text: songText,
	}

	// Настраиваем ожидания
	mockRepo.On("GetByID", "1").Return(song, nil)

	text, err := service.GetSongText("1", 1, 1)
	assert.Nil(t, err)
	assert.Equal(t, "Verse 2", text)
	mockRepo.AssertExpectations(t)
}
