package storage

import (
	"fmt"
	"ogen-film-list/gen/filmsapi"
	"sync"

	"github.com/google/uuid"
)

type FilmStorage struct {
	Mu    sync.RWMutex
	Films map[string]*filmsapi.Film
}

func NewFilmStorage() *FilmStorage {
	return &FilmStorage{
		Films: make(map[string]*filmsapi.Film),
	}
}

func (s *FilmStorage) Create(film *filmsapi.CreateFilmRequest) *filmsapi.Film {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	id := uuid.New()
	newFilm := &filmsapi.Film{
		ID:        id,
		Title:     film.Title,
		Year:      film.Year,
		Country:   film.Country,
		Director:  film.Director,
		AgeRating: filmsapi.FilmAgeRating(film.AgeRating),
		Duration:  film.Duration,
		Rating:    0,
		Actors:    []string{},
	}

	s.Films[id.String()] = newFilm
	return newFilm
}

func (s *FilmStorage) GetByID(id string) (*filmsapi.Film, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	film, exists := s.Films[id]
	return film, exists
}

func (s *FilmStorage) List(limit, offset int) []filmsapi.Film {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	films := make([]filmsapi.Film, 0)
	count := 0
	skipped := 0

	for _, film := range s.Films {
		if skipped < offset {
			skipped++
			continue
		}
		if count >= limit {
			break
		}
		films = append(films, *film)
		count++
	}

	return films
}

func (s *FilmStorage) Update(id string, update *filmsapi.UpdateFilmRequest) (*filmsapi.Film, bool) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	fmt.Printf("Storage Update - looking for ID: %s\n", id)

	film, exists := s.Films[id]
	if !exists {
		fmt.Printf("Storage Update - ID %s NOT FOUND\n", id)
		return nil, false
	}

	fmt.Printf("Storage Update - found film: %s\n", film.Title)

	film.Title = update.Title
	film.Year = update.Year
	film.Country = update.Country
	film.Rating = update.Rating
	film.Actors = update.Actors
	film.Director = update.Director
	film.AgeRating = filmsapi.FilmAgeRating(update.AgeRating)
	film.Duration = update.Duration

	return film, true
}

func (s *FilmStorage) Delete(id string) bool {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if _, exists := s.Films[id]; exists {
		delete(s.Films, id)
		return true
	}
	return false
}
