package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"ogen-film-list/gen/filmsapi"
	"ogen-film-list/internal/storage"
)

type Server struct {
	storage *storage.FilmStorage
}

var _ filmsapi.Handler = (*Server)(nil)

func NewServer() *Server {
	return &Server{
		storage: storage.NewFilmStorage(),
	}
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (s *Server) ListFilms(ctx context.Context, params filmsapi.ListFilmsParams) (filmsapi.ListFilmsRes, error) {
	limit := 20
	if params.Limit.Set && params.Limit.Value > 0 {
		limit = params.Limit.Value
	}

	offset := 0
	if params.Offset.Set && params.Offset.Value > 0 {
		offset = params.Offset.Value
	}

	films := s.storage.List(limit, offset)
	result := filmsapi.ListFilmsOKApplicationJSON(films)
	return &result, nil
}

func (s *Server) CreateFilm(ctx context.Context, req *filmsapi.CreateFilmRequest) (filmsapi.CreateFilmRes, error) {
	fmt.Println("Create film called")
	film := s.storage.Create(req)
	return film, nil
}

func (s *Server) GetFilm(ctx context.Context, params filmsapi.GetFilmParams) (filmsapi.GetFilmRes, error) {
	film, exists := s.storage.GetByID(params.ID.String())
	if !exists {
		return &filmsapi.GetFilmNotFound{}, nil
	}
	return film, nil
}

func (s *Server) UpdateFilm(ctx context.Context, req *filmsapi.UpdateFilmRequest, params filmsapi.UpdateFilmParams) (filmsapi.UpdateFilmRes, error) {
	fmt.Printf("UpdateFilm called with: %+v\n", req)
	film, updated := s.storage.Update(params.ID.String(), req)
	if !updated {
		return nil, fmt.Errorf("film not found")
	}
	return film, nil
}

func (s *Server) DeleteFilm(ctx context.Context, params filmsapi.DeleteFilmParams) (filmsapi.DeleteFilmRes, error) {
	fmt.Println("Delete film called")
	deleted := s.storage.Delete(params.ID.String())
	if !deleted {
		return &filmsapi.DeleteFilmNotFound{}, nil
	}
	return &filmsapi.DeleteFilmNoContent{}, nil
}

func (s *Server) NewError(ctx context.Context, err error) *filmsapi.DefaultErrorStatusCode {
	fmt.Printf("NewError called with: %v, type: %T\n", err, err)

	statusCode := 500
	message := "Internal server error"

	if err != nil {
		message = err.Error()
		if err.Error() == "film not found" {
			statusCode = 404
		}
	}

	errorResponse := &filmsapi.DefaultErrorStatusCode{
		StatusCode: statusCode,
		Response: filmsapi.Error{
			Code:    int32(statusCode),
			Message: message,
		},
	}

	fmt.Printf("Returning error: %+v\n", errorResponse)
	return errorResponse
}

func main() {
	server := NewServer()

	srv, err := filmsapi.NewServer(server)
	if err != nil {
		log.Fatal("Failed to create server:", err)
	}

	fmt.Println("Starting server on :8001")
	if err := http.ListenAndServe(":8001", srv); err != nil {
		log.Fatal("Server failed:", err)
	}
}
